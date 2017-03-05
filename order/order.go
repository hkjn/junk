// The order tool is a simple order matching engine.
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type (
	// OrderNumber is a unique number of an Order.
	OrderNumber int64

	// Price is a value in USD.
	Price float64

	// Volume is the amount of units being traded.
	Volume uint64

	// OrderType is the type of an order.
	OrderType uint8

	// OrderSide is the side of an order, i.e. "buy" or "sell".
	OrderSide uint8

	// Order is a request to trade items under some conditions.
	Order struct {
		// The unique id of the Order.
		id OrderNumber
		// cancelled is true if the Order has been cancelled.
		cancelled bool
		// executed is true if the Order has been fully executed.
		executed bool
		// stopTriggered is true for stop orders that have been triggered.
		stopTriggered bool
		// Type is the kind of order, which determines when the order can
		// execute.
		Type OrderType
		// Side is the direction of the order (buy/sell).
		Side OrderSide
		// Limit is the price boundary of an order, where applicable. Semantics
		// differ on OrderType.
		Limit Price
		// Volume is the number of units to trade for an order, where applicable.
		Volume Volume
		// Remaining is the number of units to trade that remains
		// unexecuted in order, where applicable.
		Remaining Volume
		// ToCancel is the OrderNumber of a previous order to cancel, where applicable.
		ToCancel OrderNumber
	}

	// orderTree is a representation of Order items for the orderbook.
	//
	// orderTree is a binary search tree that uses the price limit of orders as
	// comparison key.
	orderTree struct {
		left, right *orderTree
		item        *Order
	}

	// OrderBook holds all the orders.
	OrderBook struct {
		// buyOrders holds the BuySide orders in the book.
		buyOrders *orderTree
		// sellOrders holds the SellSide orders in the book.
		sellOrders *orderTree
		// stopOrders holds the stop orders in the book which are not yet triggered.
		stopOrders *orderTree
		nextOrder  OrderNumber
		cancelled  map[OrderNumber]bool
	}

	// Match is a match between two orders.
	Match struct {
		// Taker is the order being executed.
		Taker *Order
		// Maker is the opposing order matching Taker.
		Maker *Order
		// Volume is the number of units that orders match for.
		Volume Volume
		// Price is the price at which the orders match.
		Price Price
	}

	// Matches is several matches between pairs of orders.
	Matches []*Match
)

const (
	// Market is an Order to buy/sell at market value.
	//
	// Volume holds the number of units to trade.
	Market OrderType = iota + 1
	// Limit is an Order to buy/sell only at a specific price.
	//
	// Volume holds the number of units to trade.
	//
	// Limit holds the lower price limit to execute at for a SellSide Order.
	// Limit holds the upper price limit to execute at for a BuySide Order.
	Limit
	// Stop is an Order to trigger when price reaches given threshold.
	//
	// A Stop Order effectively creates a Market Order once the threshold
	// is reached.
	//
	// Volume holds the number of units to trade.
	//
	// Limit holds the threshold which if price goes below it triggers a SellSide Order.
	// Limit holds the threshold which if price goes above it triggers a BuySide Order.
	Stop
	// Cancel is an Order to cancel a previous Order.
	//
	// ToCancel holds the number of a previous order to cancel.
	Cancel

	// UndefinedSide is a default value for Order types which have no side.
	UndefinedSide OrderSide = iota + 1
	// BuySide is a buy-side Order.
	BuySide
	// SellSide is a sell-side Order.
	SellSide
)

var (
	orderTypes = [...]string{
		"market",
		"limit",
		"stop",
		"cancel",
	}
	orderTypesByStr = map[string]OrderType{
		"market": Market,
		"limit":  Limit,
		"stop":   Stop,
		"cancel": Cancel,
	}
)

// String returns the name of the OrderType.
func (ot OrderType) String() string { return orderTypes[ot-1] }

// String returns the name of the OrderSide.
func (oside OrderSide) String() string {
	if oside == BuySide {
		return "buy"
	} else if oside == SellSide {
		return "sell"
	} else {
		return "unknown side"
	}
}

// String returns a readable representation of the Price.
func (p Price) String() string {
	return fmt.Sprintf("%.2f", p)
}

// getMatch returns Match if the two orders can match.
func (taker *Order) getMatch(maker *Order) *Match {
	if taker.Side == maker.Side {
		// Two BuySide or SellSide orders can't possibly match.
		return nil
	}
	if taker.Type == Cancel || maker.Type == Cancel {
		// A Cancel order can't match anything.
		return nil
	}
	//if maker.Type == Stop {
	// A Stop order can't directly match anything as maker; it
	// triggers when other trades execute when its threshold is
	// reached. It can however be matched as taker.
	//return nil
	//}

	match := false
	price := Price(0.0)
	// If the taker is a Stop or Market order, any value is acceptable
	// for a match.
	//if taker.Type == Stop {
	// Stop orders need to be triggered first.
	//if taker.stopTriggered {
	//match = true
	// TODO: Correct? What if other side is Market too, with no Limit?
	//price = maker.Limit
	//}
	if taker.Type == Market {
		match = true
		price = maker.Limit
	} else if taker.Side == BuySide && taker.Limit >= maker.Limit {
		match = true
		price = taker.Limit
	} else if taker.Side == SellSide && taker.Limit <= maker.Limit {
		match = true
		price = maker.Limit
	}

	if !match {
		log.Fatalf("No match; bug?\n")
		return nil
	}

	volume := taker.Volume
	if volume > maker.Volume {
		volume = maker.Volume
	}
	debug("[MATCH] %q matches %q at %v, for %v units\n", taker, maker, price, volume)
	return &Match{
		Taker:  taker,
		Maker:  maker,
		Volume: Volume(volume),
		Price:  price,
	}
}

// String returns a description of the Order.
func (order Order) String() string {
	cancelled := ""
	if order.cancelled {
		cancelled = "[cancelled] "
	}
	executed := ""
	if order.executed {
		executed = "[executed] "
	}
	triggered := ""
	if order.stopTriggered {
		triggered = "[triggered] "
	}
	idstr := ""
	if order.id > 0 {
		idstr = fmt.Sprintf("[id %d] ", order.id)
	}
	status := fmt.Sprintf("%s%s%s%s", idstr, cancelled, executed, triggered)

	cond := "?!?"
	if order.Side == BuySide {
		cond = "<="
	} else if order.Side == SellSide {
		cond = ">="
	}

	if order.Type == Market {
		return fmt.Sprintf(
			"%s%v order to %v %v units at market price, with %v remaining",
			status,
			order.Type,
			order.Side,
			order.Volume,
			order.Remaining,
		)
	}

	if order.Type == Limit {
		return fmt.Sprintf(
			"%s%v order to %v %v units %s $%v, with %v remaining",
			status,
			order.Type,
			order.Side,
			order.Volume,
			cond,
			order.Limit,
			order.Remaining,
		)
	}

	if order.Type == Stop {
		cond := "?!?"
		if order.Side == BuySide {
			cond = ">"
		} else if order.Side == SellSide {
			cond = "<"
		}
		return fmt.Sprintf(
			"%s%v order to %v %v units if price goes %s %v, with %v remaining",
			status,
			order.Type,
			order.Side,
			order.Volume,
			cond,
			order.Limit,
			order.Remaining,
		)
	}

	if order.Type == Cancel {
		return fmt.Sprintf(
			"%s%v order that disables #%v",
			status,
			order.Type,
			order.ToCancel,
		)
	}

	return "Order{???}"
}

// String returns a readable description of the Match.
func (m Match) String() string {
	return fmt.Sprintf(
		"match between order %v and %v for %v units at $%v",
		m.Taker,
		m.Maker,
		m.Volume,
		m.Price,
	)
}

func (t orderTree) String() string {
	if t.item == nil {
		return ""
	}
	return fmt.Sprintf(" %q", t.item)
}

// Output returns the output format to emit for the Match.
func (m Match) Output() string {
	return fmt.Sprintf("match %v %v %v %v", m.Taker.id, m.Maker.id, m.Volume, m.Price)
}

// Len, Swap and Less implements sort.Sort interface for Matches.
func (ms Matches) Len() int      { return len(ms) }
func (ms Matches) Swap(i, j int) { ms[i], ms[j] = ms[j], ms[i] }

// Less returns whether the Matches with index i comes before index j.
//
// Price is checked first. Where it differs the match with best price
// for Taker comes first.
//
// If first match has lower price and first match was a sell
// order, it's better.
//
// If Price is the same, oldest Maker comes before.
func (ms Matches) Less(i, j int) bool {
	if ms[i].Price < ms[j].Price {
		return ms[i].Taker.Side == SellSide
	} else if ms[i].Price > ms[j].Price {
		return ms[i].Taker.Side != SellSide
	}
	return ms[i].Maker.id < ms[j].Maker.id
}

func info(format string, a ...interface{}) {
	if true {
		fmt.Printf("[I] "+format, a...)
	}
}

func debug(format string, a ...interface{}) {
	if true {
		fmt.Printf("[D]   "+format, a...)
	}
}

func debugv(format string, a ...interface{}) {
	if false {
		fmt.Printf("[DD]     "+format, a...)
	}
}

// newOrder returns the Order parsed from orderstr.
//
// newOrder panics if orderstr is invalid.
func newOrder(orderstr string) *Order {
	parts := strings.Split(orderstr, " ")
	if len(parts) != 4 {
		log.Fatalf("Unexpected order string: %q\n", orderstr)
	}

	order := &Order{}
	ot, ok := orderTypesByStr[parts[0]]
	if !ok {
		log.Fatalf("Unexpected order type: %q\n", parts[0])
	}
	order.Type = ot

	if order.Type != Cancel {
		if parts[1] == "buy" {
			order.Side = BuySide
		} else if parts[1] == "sell" {
			order.Side = SellSide
		}
	}

	value1, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		log.Fatalf("Unexpected value1: %q\n", parts[2])
	}
	if order.Type == Cancel {
		order.ToCancel = OrderNumber(value1)
	} else {
		order.Volume = Volume(value1)
		order.Remaining = order.Volume
	}

	value2, err := strconv.ParseFloat(parts[3], 64)
	if err != nil {
		log.Fatalf("Unexpected value2: %q\n", parts[3])
	}
	if order.Type == Limit || order.Type == Stop {
		order.Limit = Price(value2)
	}

	return order
}

// newOrderBook returns a new OrderBook.
func newOrderBook() OrderBook {
	return OrderBook{
		nextOrder: 1,
		cancelled: map[OrderNumber]bool{},
	}
}

// insert adds a new Order to the orders BST.
//func (t *orderTree) insert(order *Order) {
//	debug("insert(%q)\n", order)
//	c := t
//	for c != nil {
//		if order.Limit <= c.item.Limit {
//			debugv("inserting to left (lower)\n")
//			c.left = &orderTree{}
//			c = c.left
//		} else {
//			debugv("inserting to right (higher)\n")
//			c.right = &orderTree{}
//			c = c.right
//		}
//	}
//	// TODO: Handle missing Limit; c.key should be nil
//	// c.key = order.id
//	c.item = order
//}

// insert adds a new order to the orderTree and returns the new tree.
//
// TODO: Change from recursion to iteration, make method on orderTree.
func insert(node *orderTree, order *Order) *orderTree {
	if node == nil {
		node = &orderTree{
			item: order,
		}
	} else if order.Limit <= node.item.Limit {
		debugv("inserting %v to left (lower)\n", order.Limit)
		node.left = insert(node.left, order)
	} else {
		debugv("inserting %v to right (higher)\n", order.Limit)
		node.right = insert(node.right, order)
	}
	return node
}

func traverse(t *orderTree, desc string) {
	if t == nil {
		return
	}
	traverse(t.left, desc+"L")
	debug("traversing [%v] %v: %+v\n", desc, t.item.id, t.item)
	traverse(t.right, desc+"R")
}

// exec traverses the BST and returns the Match, if any.
//
// If max is true, orders with limit values >= specified value are returned.
// If max is false, orders with limit values <= specified value are returned.
func (node *orderTree) exec(taker *Order, wantHighest bool) *Match {
	c := node
	makers := []*Order{}
	debug("exec(%q)\n", taker)
	// TODO: Handle missing Limit for Market
	for c != nil {
		if wantHighest {
			if c.item.Limit >= taker.Limit {
				makers = append(makers, c.item)
				debugv("fetched order: %v\n", makers)
			}
			debugv("fetching right (higher)\n")
			c = c.right
		} else {
			if c.item.Limit <= taker.Limit {
				makers = append(makers, c.item)
				debugv("fetched order: %v\n", makers)
			}
			debugv("fetching left (lower)\n")
			c = c.left
		}
		time.Sleep(time.Millisecond * 200) // TODO: Remove
	}

	matches := Matches{}
	for _, maker := range makers {
		matches = append(matches, taker.getMatch(maker))
	}
	sort.Sort(sort.Reverse(matches))

	if len(matches) == 0 {
		debug("exec() found no matches\n")
		return nil
	}
	debugv("exec() got matches: %v\n", matches)

	// Execute the order.
	match := matches[0]
	match.Maker.Volume -= match.Volume
	if match.Maker.Volume <= 0.0 {
		match.Maker.executed = true
	}
	match.Taker.Volume -= match.Volume
	if match.Taker.Volume <= 0.0 {
		match.Taker.executed = true
	}
	debug("exec() excuted %q\n", match)
	return match
}

// delete removes an Order from the BST.
func (t *orderTree) delete(order *Order) {
	// TODO: Need a reference to the tree here.. even if we store a pointer
	// back to the *orderTree, we'd still need a reference to its parent to
	// be able to drop that parent's reference when we delete
	// though. Looking up node by order.id is O(n).
	debug("delete(%v)\n", order)
}

// findStops returns all Stop orders triggered by specified price.
func (book *OrderBook) findStops(order *Order) []*Order {
	// TODO: Traverse stopOrders to find if order triggers any stop
	// order. If it did, we should add that stop order to buyOrders or
	// sellOrders, and execute that order if so.

	c := book.stopOrders
	wantHighest := order.Side == BuySide
	triggered := []*Order{}
	debug("findStops(%v)\n", order)
	for c != nil {
		if wantHighest {
			if c.item.Limit >= order.Limit {
				triggered = append(triggered, c.item)
				debugv("fetched order: %v\n", triggered)
			}
			debugv("fetching right (higher)\n")
			c = c.right
		} else {
			if c.item.Limit <= order.Limit {
				triggered = append(triggered, c.item)
				debugv("fetched order: %v\n", triggered)
			}
			debugv("fetching left (lower)\n")
			c = c.left
		}
		time.Sleep(time.Millisecond * 200) // TODO: Remove
	}

	return triggered
}

// Add adds and attempts to execute an Order.
//
// The new Order are matched with existing orders in the book. Order
// matching depends on the type.
//
// If there's matching orders, they are executed, and the resulting
// matches are returned.
func (book *OrderBook) Add(taker *Order) Matches {
	taker.id = book.nextOrder
	book.nextOrder++
	if taker.Type == Cancel {
		book.cancelled[taker.ToCancel] = true
		// TODO: Delete taker.ToCancel from appropriate tree here by
		// traversing them and then deleting the order.
		debug("Cancelled %v\n", taker.ToCancel)
		return nil
	}
	if taker.Type == Stop {
		debug("Added stop order %v\n", taker)
		book.stopOrders = insert(book.stopOrders, taker)
		return nil
	}

	if taker.Side == BuySide {
		book.buyOrders = insert(book.buyOrders, taker)
	} else {
		book.sellOrders = insert(book.sellOrders, taker)
	}
	info("Added order %q\n", taker)

	// TODO: Need to pass on a reference to the book.buyOrders *orderTree here to exec, so it can

	// Look for matches for the recently added order.
	matches := Matches{}
	for !taker.executed {
		var match *Match
		if taker.Side == BuySide {
			match = book.sellOrders.exec(taker, false)
		} else {
			match = book.buyOrders.exec(taker, true)
		}
		if match == nil {
			debug("No new matches, returning the ones we have: %v\n", matches)
			return matches
		}
		debug("Adding best new match found: %v\n", match)
		matches = append(matches, match)
		// TODO: Both taker and maker are executed if we made it here, need
		// to be deleted from BST.
		if taker.Side == BuySide {
			if match.Taker.executed {
				book.buyOrders.delete(match.Taker)
			}
			if match.Maker.executed {
				book.sellOrders.delete(match.Maker)
			}
		} else {
			if match.Taker.executed {
				book.sellOrders.delete(match.Taker)
			}
			if match.Maker.executed {
				book.buyOrders.delete(match.Maker)
			}
		}
	}
	return matches
}

// getStopMatches returns any matches for stop orders triggered by the matches.
func (book *OrderBook) getTriggeredStops(oldMatches Matches) Matches {
	matches := Matches{}
	for _, oldMatch := range oldMatches {
		// TODO: Each oldMatch.Taker here is an executed order. We know
		// the oldMatch.Price, and now need to traverse the BST to find if
		// this price is below any SellSide Stop order's Limit, or above
		// any BuySide Stop order's Limit. If so, they trigger and we want
		// to match + execute them, and return the matches.

		for _, taker := range book.findStops(oldMatch.Taker) {
			// Any Order in triggered is a Stop order which was triggered by
			// recent executions. They should be removed from stopOrders, and
			// added to buyOrders / sellOrders.
			debug("Should delete stopOrder %v\n", taker)
			book.stopOrders.delete(taker)
			if taker.Side == BuySide {
				debug("Adding triggered stoporder to buyOrders: %q\n", taker)
				book.buyOrders = insert(book.buyOrders, taker)
			} else {
				debug("Adding triggered stoporder to sellOrders: %q\n", taker)
				book.sellOrders = insert(book.sellOrders, taker)
			}

			// Look for matches for the recently triggered stop order.
			for !taker.executed {
				var match *Match
				if taker.Side == BuySide {
					match = book.sellOrders.exec(taker, false)
				} else {
					match = book.buyOrders.exec(taker, true)
				}
				if match == nil {
					debug("No new matches, returning the ones we have: %v\n", matches)
					return matches
				}
				debug("Adding best new match found: %v\n", match)
				matches = append(matches, match)
			}
		}
	}
	return matches
}

func main() {
	book := newOrderBook()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		orderstr := scanner.Text()
		if err := scanner.Err(); err != nil {
			log.Fatalf("Failed to read standard input: %v\n", err)
		}
		order := newOrder(orderstr)
		// TODO: Maybe type returned here should be Executions or
		// something; "matches" is misleading since we already executed them.
		matches := book.Add(order)
		for _, match := range matches {
			fmt.Println(match.Output())
		}
		// The matches for order might have triggered some stop orders.
		for _, match := range book.getTriggeredStops(matches) {
			fmt.Println(match.Output())
		}
	}
}
