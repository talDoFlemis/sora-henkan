import http from "k6/http";
import { check, sleep } from "k6";
import { SharedArray } from "k6/data";
import animes_urls from "./animes.json";

// --- Configuration ---

// Use the __ENV object to read the TARGET_ENDPOINT from the environment
const DOMAIN = __ENV.TARGET_ENDPOINT ?? "https://api-sorahenkan.flemis.cloud";

// Fallback or validation: Check if the environment variable was provided
if (!DOMAIN) {
  throw new Error(
    "TARGET_ENDPOINT environment variable is not set. Please set it before running.",
  );
}

const TARGET_ENDPOINT = `${DOMAIN}/v1/images/`;

// Set a range for random width and height (e.g., from 100 to 1000 pixels)
const MIN_DIMENSION = 100;
const MAX_DIMENSION = 5000;

// k6 options (Load configuration)
export const options = {
  // A sample load test profile: 10 Virtual Users for 30 seconds
  stages: [
    { duration: "10s", target: 5 }, // ramp up to 5 VUs
    { duration: "15s", target: 10 }, // stay at 10 VUs
    { duration: "5s", target: 0 }, // ramp down to 0 VUs
  ],
  thresholds: {
    // Fail the test if 95% of requests take longer than 500ms
    http_req_duration: ["p(95)<500"],
    // Fail the test if the request failure rate is above 1%
    checks: ["rate>0.99"],
  },
};

// --- Data Loading (Runs only once before the test starts) ---

// --- Main Execution Function (Runs repeatedly by each Virtual User) ---

export default function () {
  // 1. Select a random image URL
  const randomUrlIndex = Math.floor(Math.random() * animes_urls.length);
  const imageUrl = animes_urls[randomUrlIndex];

  // 2. Generate random width and height
  const randomWidth =
    Math.floor(Math.random() * (MAX_DIMENSION - MIN_DIMENSION + 1)) +
    MIN_DIMENSION;
  const randomHeight =
    Math.floor(Math.random() * (MAX_DIMENSION - MIN_DIMENSION + 1)) +
    MIN_DIMENSION;

  // 3. Construct the JSON payload
  const payload = JSON.stringify({
    image_url: imageUrl,
    scale: {
      enabled: true,
      width: randomWidth,
      height: randomHeight,
    },
  });

  // Headers to specify content type
  const params = {
    headers: {
      "Content-Type": "application/json",
    },
  };

  // 4. Send the POST request (simulating the cURL call)
  // The endpoint is dynamically read from the environment variable
  const res = http.post(TARGET_ENDPOINT, payload, params);

  // 5. Check and sleep
  check(res, {
    "is status 200": (r) => r.status === 200,
    "has body": (r) => r.body.length > 0,
  });

  // Control the request rate - e.g., wait 0.5 seconds between iterations
  sleep(0.5);
}
