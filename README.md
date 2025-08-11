# GO Rate Limiter: A Customizable, Standalone Rate Limiting  Service
A fully customizable Rate Limiting service written in Go. With a simple `JSON` configuration file, you can define per-route rate limits, choose the algorithm, and set the request interval, all without modifying any code. 

Perfect for education, testing, or production use, this project lets you experiment with different rate limiting algorithms, add logging to measure performance, and choose the right strategy for each endpoint on your own. Being written in Go means it's secure, efficient, and highly concurrent, making it reliable to use for preventing outages and abuse.
# Table of Contents
1. [Features](#features)
	- [Custom Configurations](#custom-configurations)
	- [Rate Limiting Key Variety](#rate-limiting-key-variety)
	- [Rate Limiting Algorithms Variety](#rate-limiting-algorithms-variety)
	- [Fast Processing Time](#fast-processing-time)
2. [Rate Limiting Algorithms](#rate-limiting-algorithms)
3. [How to Use](#how-to-use)
	- [Why standalone?](#why-standalone)
	- [Usage Steps](#usage-steps)
4. [Final Node.js Server Test](#final-nodejs-server-test)
5. [Extra Words](#extra-words)


# Features

### Custom Configurations
You can create fully customized configurations, select algorithm to use (see [Rate Limiting Algorithms](#Rate\Limiting\Algorithms)), set the time interval, and define maximum number of requests in that interval for each route.
### Rate Limiting Key Variety
You can choose what client key to rate-limit based on, by setting the request header `X-Rate-Limit-Key` to your desired option. Most common options are a user's IP, API Key, or his user id if he's an authenticated user. You can choose any identifier you prefer.
### Rate Limiting Algorithms Variety
Many libraries only offer one default algorithm for all routes, but this approach isn’t always optimal. And according to what I've learned, which is **rarely** optimal. Each algorithm has strengths and weaknesses, performs best on specific type of routes while underperforming on others.
### Fast Processing Time
I have utilized `goroutines` and production-ready code, in addition to the nature of GO, that resulted in this service to be fast and processes requests in no time. Benchmarks show an average overhead of just 0-1 ms per request, with rare peaks of 2-3 ms.
# Rate Limiting Algorithms
Here's a list of used algorithms that you can use by inserting its corresponding name in your `json` configurations file. You can check [this video](https://youtu.be/mQCJJqUfn9Y) to learn more about each algorithm, their pros and cons, and their best use cases.
- `FixedWindowCounter`: The most straight-forward algorithm that is used by default in `express-rate-limiting` library in `node.js`.
- `SlidingWindowLog`: Another variation of utilizing windows in rate limiting. 
- `SlidingWindowCounter`: A very interesting approach that takes the best of both previous windows variations, it's less intuitive but it is so much worth investigating into.
- `TokenBucket`: An algorithm that doesn't utilize windows, yet focuses on a fixed rate per unit time, to eventually achieve the desired limiting interval.
Many other algorithms exist, but these are the most widely used based on my research.
# Configuration
Provided is an example for a simple `config.json` file, that limits the login route requests to 5 requests in a minute, and the posts route requests to 10 requests per 2 minutes maximum.
Example `config.json`:
```json
{
  "algorithm": "SlidingWindowCounter",
  "routes": [
    {
      "route": "/api/v1/login",
      "limit": 5,
      "interval": "1m"
    },
    {
      "route": "/api/v1/posts",
      "limit": 10,
      "interval": "2m"
    }
  ]
}
```
The config file can be placed anywhere, and you can have multiple config files for multiple instances. The config file is passed as an argument when starting the application, default is `./config/config.json`.
# Installation
No need to clone the repo, the image is pushed to dockerhub.
1. Have docker installed
2. Create this `docker-compose.yml` file anywhere
```yaml
# Example docker-compose.yml
services:
rate-limiter:
  image: loayahmed/go-rate-limiter:latest
  ports:
    - "4000:9240"
  volumes:
    - ./configs:/home/app/configs
  command: ["./rate-limiter", "-f=./configs/config.json", "-p=:9240"]
```
> Change your config file directory and update it in the CLI arguments

That's it, congratulations, you have it running now!
# How to Use
I have decided to make this Rate Limiting project as a standalone service by design. 
### Why standalone?
- To make it work with any backend framework
- Offloads computation from your main backend.
### Usage Steps
It is used as a simple backend API request, you can send any type of request (preferred to be a GET or POST request) to the IP and port of the running instance. You should create a middleware function that sends a request to the rate limiter. Your middleware will typically do the follows:
1. Set a header of `X-Original-Path` to your request URL (like `/api/v1/login`).
	- If the header is not found, the rate limiter will use its own requested URL path.
2. Set the header of `X-Rate-Limit-Key` to your desired key to rate limit based on. Note that this header **is required** and will not be defaulted by the rate limiter, since the client's original IP can't be accessed by the rate limiter.
3. Send the request and expect to receive a status code of `429` indicating too many requests, in addition to the time remaining in the body, or a `200` code.
#### Below is an example middleware in Node.js that checks the rate limiter before processing requests.
```js
// example server.js
const rateLimit = async (req, res, next) => {
	try {
	const response = await fetch("http://localhost:4000", {
		headers: {
		// or use userId or whatever key identifier you prefer
		"X-Rate-Limit-Key": req.ip, 
		
		// the path must be included in the rate limiter config.json,
		// otherwise it'll be ignored (always passed)
		"X-Original-Path": req.path, 
	},
	});
	if (response.status === 429) {
		const timeRemaining = response.headers.get("Retry-After"); // set by the rate limiter,
		// you can choose to either display it or not
		const message = "Too many requests.";
		return res.status(429).send(message + " Retry after: " + timeRemaining);
	}
	next();
	} catch (err) {
		console.error("Error fetching rate limit:", err);
		return res.status(500).send("Internal Server Error");
	}
};
```

> Also fetching for the path itself will work, if the header is not set.
> For example, this will also work:

```js
// example inside rateLimit function
fetch(`http://localhost:4000/${req.path}`, {
	headers: {
		"X-Rate-Limit-Key": req.ip,
	}
})
```

Now you can inject your middleware anywhere you want, or you can even use it for every single request and the non-configured routes will be passed (although this is not recommended). Here's an example for such usage:
```js
// Example server.js
app.get("/api/v1/login", rateLimit, handleLogin);
app.get("/api/v1/posts", rateLimit, handlePosts);
```

```js
// Example server.js

// not recommended, but still valid
app.use(rateLimit);
app.get("/api/v1/login", handleLogin);
app.get("/api/v1/posts", handlePosts);
```


# Final Node.js Server Test
```js
// Example index.js
const express = require("express");

const app = express();

// set the request start time
app.use((req, res, next) => {
  res.locals.startTime = Date.now();
  next();
});

const rateLimit = async (req, res, next) => {
  try {
    fetch("http://localhost:4000", {
      headers: {
        "X-Rate-Limit-Key": req.ip, // or use userId or whatever key identifier you prefer
        "X-Original-Path": req.path, // the path must be included in the rate limiter config.json,
        // otherwise it'll be ignored (always passed)
      },
    }).then((resp) => {
      if (resp.status === 429) {
        const timeRemaining = resp.headers.get("Retry-After"); // set by the rate limiter,
        //  you can choose to either display it or not

        const message = "Too many requests.";

        return res.status(429).send(message + " Retry after: " + timeRemaining);
      }
      next();
    });
  } catch (err) {
    console.error("Error fetching rate limit:", err);
    return res.status(500).send("Internal Server Error");
  }
};

app.get("/api/v1/posts", rateLimit, (_, res) => {
  res.send(
    "Posts List, time taken since receiving the request: " +
      (Date.now() - res.locals.startTime) +
      "ms"
  );
});
app.get("/api/v1/login", rateLimit, (_, res) => {
  res.send(
    "Login Page, time taken since receiving the request: " +
      (Date.now() - res.locals.startTime) +
      "ms"
  );
});

app.listen(3000, () => console.log("Server is running on port 3000"));
```
# Extra Words
 I built this project to deepen my Go skills and create a robust, production-ready service. I’m proud of the result and welcome feedback or contributions. I’ve also read extensively about writing idiomatic Go code to ensure the architecture and implementation follow best Go practices.

**Thank you.**
