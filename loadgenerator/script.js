import http from "k6/http";
import { check, sleep } from "k6";

const animes_urls = [
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/about-to-die.jpg",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/access_control.jpg",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/aegis-pointing-upclose.jpg",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/aegis-pointing.png",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/aegis-pointing2.png",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/alexjones.png",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/aoc.jpg",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/archie.png",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/bang-bang-white.png",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/beeeeginbot.jpeg",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/black-hole-up-close.png",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/blue-chains.png",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/blue-orange.jpeg",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/blue-red-blue-blond.png",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/blue_lookback.png",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/both-popes.png",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/Bournane_Abdelkrim_the_primeagen_3400abb7-3a41-4824-9582-5cae521b42b4.png",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/butterfly.jpg",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/car_girl_white_dress.jpeg",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/catumbrella.jpg",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/changing2.png",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/charging.png",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/charging3.png",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/city-scape.gif",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/crown-blackhole.png",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/crown-close-up.png",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/crying_vscode.jpg",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/cyber_girl.jpg",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/dead_eyes.jpg",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/Diversity.png",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/doingnumbers.jpeg",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/dragon-girl.jpg",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/driving.jpg",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/endless-summer.jpg",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/garden.jpeg",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/girl-gun-walk-away.jpg",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/girl-gun.png",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/girl-with-gu.png",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/green-brown-yellow.jpg",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/green-eyes.jpg",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/green-green-red.jpg",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/grung-green-yellow-refd.png",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/hatsune_running.jpg",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/in-front-of-car.jpg",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/in-the-rain.png",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/is_this_zelda.png",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/javascript.jpg",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/lose-control.png",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/moon-upside.jpg",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/more_guns_walking.jpg",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/mythra.png",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/n8-versace.png",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/nicerocket.jpeg",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/nicerocket.png",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/not-stonks.png",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/nothing_but_couch.jpg",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/nr-b2.jpg",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/omao_1.jpeg",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/ONLY_WEEBS.jpeg",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/ONLY_WEEBS_2.jpeg",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/onlyfans.png",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/oskr_the_primeagen_6371be34-bd8a-4643-82c1-d480ec36ea29.png",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/powder.jpg",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/prime-vs-erica.png",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/primebffs.png",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/purple-sword.jpg",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/purple-two-swords.jpg",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/purple.jpg",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/raid-in-the-dark.png",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/README.md",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/RETF.jpg",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/rose-pine-sword.jpg",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/rust-is-bad.jpg",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/rust.jpg",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/rustprogrammers.jpg",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/SecretofMana.jpg",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/Selection_259.png",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/shestoocool.png",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/skyline.jpg",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/slapping.jpg",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/snapn_cards.jpg",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/sword-red.png",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/ThePrimeagen.jpeg",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/thisguyvims.png",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/umbrella.jpeg",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/uwuntu.png",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/white-red-green.jpg",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/yayayayayaya.jpg",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/yellow_red_blue.png",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/zelda-white.bmp",
  "https://raw.githubusercontent.com/ThePrimeagen/anime/refs/heads/master/zelda_look_alike.svg",
];

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
