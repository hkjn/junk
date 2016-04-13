package hello;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class Greeting {
	static final Logger LOG = LoggerFactory.getLogger(Greeting.class);
	private final long id;
	private final String content;

	public Greeting(long id, String content) {
		System.out.println("FIXMEH: Greeting()");
		LOG.info("FIXMEH: info() logging says {} from {}", "hi", this);
		LOG.error("FIXMEH: error()");
		LOG.debug("FIXMEH: debug()");
		this.id = id;
		this.content = content;
	}

	public long getId() {
		return id;
	}

	public String getContent() {
		return content;
	}
}
