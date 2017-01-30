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
)

type (
	// OrderNumber is a unique number of an Order.
	OrderNumber int64

	// Price is a value in some currency.
	Price float64

	// Volume is the amount of units being traded.
	Volume uint64

	// OrderType is the type of an order.
	OrderType uint8

	// OrderSide is the side of an order, i.e. "buy" or "sell".
	OrderSide uint8

	// Value1 represents either the number of units to trade or the
	// order number to cancel, depending on OrderType.
	Value1 uint64

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
		// Type is the kind of OrderType, which determines when the order
		// can execute.
		Type OrderType
		// Side is the direction of the order (buy/sell).
		Side OrderSide
		// Value1 is the first value of an order, with semantics differing
		// by OrderType.
		Value1 Value1
		// Value2 is the second value of an order, with semantics differing
		// by OrderType.
		Value2 Price
		// Volume is the number of units to trade that remains unexecuted
		// in order, where applicable.
		Volume Volume
	}

	// OrderBook holds all the orders.
	OrderBook struct {
		// TODO: Might be better to keep index of price bands to orders,
		// or at least to not scan through orders we already have
		// executed. Maybe https://en.wikipedia.org/wiki/Skip_list? Though
		// https://en.wikipedia.org/wiki/Binary_search_tree is likely
		// simplest. Or red-black tree?
		orders     map[OrderNumber]Order
		stopOrders map[OrderNumber]Order
		nextOrder  OrderNumber
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
	// Value1 is the number of units to trade.
	//
	// Value2 is ignored.
	Market OrderType = iota + 1
	// Limit is an Order to buy/sell only at a specific price.
	//
	// Value1 is the number of units to trade.
	//
	// Value2 is the lower price limit to execute at for a SellSide Order.
	// Value2 is the upper price limit to execute at for a BuySide Order.
	Limit
	// Stop is an Order to trigger when price reaches given threshold.
	//
	// A Stop Order effectively creates a Market Order once the threshold
	// is reached.
	//
	// Value1 is the number of units to trade.
	//
	// Value2 is the threshold which if price goes below it triggers a SellSide Order.
	// Value2 is the threshold which if price goes above it triggers a BuySide Order.
	Stop
	// Cancel is an Order to cancel a previous Order.
	//
	// Value1 is the number of a previous Order to cancel.
	//
	// Value2 is ignored.
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
func (taker Order) getMatch(maker Order) *Match {
	if taker.Side == maker.Side {
		// Two BuySide or SellSide orders can't possibly match.
		return nil
	}
	if taker.Type == Cancel || maker.Type == Cancel {
		// A Cancel order can't match anything.
		return nil
	}
	if maker.Type == Stop {
		// A Stop order can't directly match anything as maker; it
		// triggers when other trades execute when its threshold is
		// reached. It can however be matched as taker.
		return nil
	}

	match := false
	price := Price(0.0)
	// If the taker is a Stop or Market order, any value is acceptable
	// for a match.
	if taker.Type == Stop {
		if taker.stopTriggered {
			match = true
			price = maker.Value2
		}
	} else if taker.Type == Market {
		match = true
		price = maker.Value2
	} else if taker.Side == BuySide && taker.Value2 >= maker.Value2 {
		match = true
		price = taker.Value2
	} else if taker.Side == SellSide && taker.Value2 <= maker.Value2 {
		match = true
		price = maker.Value2
	}

	if !match {
		return nil
	}

	volume := taker.Volume
	if volume > maker.Volume {
		volume = maker.Volume
	}
	debug("%q matches %q at %v, for %v units\n", taker, maker, price, volume)
	return &Match{
		Taker:  &taker,
		Maker:  &maker,
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
	idstr := ""
	if order.id > 0 {
		idstr = fmt.Sprintf("[id %d] ", order.id)
	}
	desc := fmt.Sprintf("%s%s%s", idstr, cancelled, executed)
	if order.Type == Cancel {
		return fmt.Sprintf(
			"%scancel order for #%v",
			desc,
			order.Value1,
		)
	}
	if order.Type == Stop {
		if order.Side == BuySide {
			return fmt.Sprintf(
				"%sstop buy order for #%v which triggers if price goes > $%v, with %v remaining",
				desc,
				order.Value1,
				order.Value2,
				order.Volume,
			)
		} else {
			return fmt.Sprintf(
				"%sstop sell order for #%v which triggers if price goes < $%v, with %v remaining",
				desc,
				order.Value1,
				order.Value2,
				order.Volume,
			)
		}
	}
	cond := "?!?"
	if order.Side == BuySide {
		cond = "<="
	} else if order.Side == SellSide {
		cond = ">="
	}
	if order.Type == Market {
		return fmt.Sprintf(
			"%smarket order to %v %v units %s market price, with %v remaining",
			desc,
			order.Side,
			order.Value1,
			cond,
			order.Volume,
		)
	}
	return fmt.Sprintf(
		"%s%v order to %v %v units at %s $%v, with %v remaining",
		desc,
		order.Type,
		order.Side,
		order.Value1,
		cond,
		order.Value2,
		order.Volume,
	)
}

// String returns a readable description of the Match.
func (m Match) String() string {
	return fmt.Sprintf(
		"match between order %q and %q for %v units at $%v",
		m.Taker,
		m.Maker,
		m.Volume,
		m.Price,
	)
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

func debug(format string, a ...interface{}) {
	if false {
		fmt.Printf("[DEBUG] "+format, a...)
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

	otype, ok := orderTypesByStr[parts[0]]
	if !ok {
		log.Fatalf("Unexpected order type: %q\n", parts[0])
	}

	var oside OrderSide
	if otype != Cancel {
		if parts[1] == "buy" {
			oside = BuySide
		} else if parts[1] == "sell" {
			oside = SellSide
		} else {
			log.Fatalf("Unexpected order side: %q\n", parts[1])
		}
	}

	value1, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		log.Fatalf("Unexpected value1: %q\n", parts[2])
	}
	value2, err := strconv.ParseFloat(parts[3], 64)
	if err != nil {
		log.Fatalf("Unexpected value2: %q\n", parts[3])
	}

	return &Order{
		Type:   otype,
		Side:   oside,
		Value1: Value1(value1),
		Value2: Price(value2),
		Volume: Volume(value1),
	}
}

// newOrderBook returns a new OrderBook.
func newOrderBook() OrderBook {
	return OrderBook{
		orders:     map[OrderNumber]Order{},
		stopOrders: map[OrderNumber]Order{},
		nextOrder:  1,
	}
}

// getMatches returns all Matches to specified Order.
//
// If multiple orders match the new order, price limit is the first
// priority.
//
// If several matching orders have the same price limit, oldest order
// is matched first.
func (book *OrderBook) getMatches(taker *Order) Matches {
	matches := Matches{}
	for _, maker := range book.orders {
		if maker.cancelled || maker.executed {
			continue
		}
		match := taker.getMatch(maker)
		if match != nil {
			matches = append(matches, match)
		}
	}
	sort.Sort(sort.Reverse(matches))
	return matches
}

// Add adds and attempts to execute an Order.
//
// The new Order are matched with existing orders in the book. Order
// matching depends on the type.
//
// If there's matching orders, they are executed, and the resulting
// matches are returned.
func (book *OrderBook) Add(order *Order) Matches {
	order.id = book.nextOrder
	book.orders[order.id] = *order
	if order.Type == Stop {
		book.stopOrders[order.id] = *order
	}
	book.nextOrder++
	debug("Added order %q\n", order)

	if order.Type == Cancel {
		toCancel, exists := book.orders[OrderNumber(order.Value1)]
		if exists && !toCancel.executed {
			// If order to cancel didn't exist or was fully executed, it
			// can't be cancelled. Weird, but this is supposed to be no-op.
			toCancel.cancelled = true
			book.orders[toCancel.id] = toCancel
			debug("Cancelled %q\n", toCancel)
			_, exists := book.stopOrders[toCancel.id]
			if exists {
				book.stopOrders[toCancel.id] = toCancel
			}
		}
		return nil
	}

	matches := Matches{}
	// Stop orders don't match as taker, but they can be triggered
	// after other orders execute.
	for order.Type != Stop && !order.executed {
		// Look through existing orders in book for ones matching price
		// limit new order wants, if so they match and can be executed.
		newMatches := book.getMatches(order)
		if len(newMatches) == 0 {
			return matches
		}
		order = book.execute(newMatches[0])
		matches = append(matches, newMatches[0])
	}
	return matches
}

func (book *OrderBook) execute(match *Match) *Order {
	debug("Executing: %v\n", match)
	match.Maker.Volume -= match.Volume
	if match.Maker.Volume <= 0.0 {
		match.Maker.executed = true
		book.orders[match.Maker.id] = *match.Maker
	}
	match.Taker.Volume -= match.Volume
	if match.Taker.Volume <= 0.0 {
		match.Taker.executed = true
		book.orders[match.Taker.id] = *match.Taker
	}
	return match.Taker
}

// getStopMatches returns any matches for stop orders triggered by the match.
func (book *OrderBook) getTriggeredStops(oldMatches Matches) Matches {
	matches := Matches{}
	for _, oldMatch := range oldMatches {
		for _, stop := range book.stopOrders {
			stopp := &stop
			if stop.Side == BuySide && oldMatch.Price > stop.Value2 {
				// The match triggers this BuySide Stop.
				stopp.stopTriggered = true
			} else if stop.Side == SellSide && oldMatch.Price < stop.Value2 {
				// The match triggers this SellSide Stop.
				stopp.stopTriggered = true
			}
			for !stopp.executed {
				newMatches := book.getMatches(stopp)
				if len(newMatches) == 0 {
					break
				} else {
					book.execute(newMatches[0])
					stopp = newMatches[0].Taker
					book.stopOrders[stopp.id] = *stopp
					debug("stop order was executed: %v\n", stopp)
					matches = append(matches, newMatches[0])
				}
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
		matches := book.Add(order)
		for _, match := range matches {
			fmt.Println(match.Output())
		}
		// The matches for order might also trigger stop orders.
		for _, match := range book.getTriggeredStops(matches) {
			fmt.Println(match.Output())
		}
	}
}
