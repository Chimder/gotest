import http from "k6/http";
import { check, sleep } from "k6";
import { Counter } from "k6/metrics";

export const requests = new Counter("http_reqs");

export const options = {
  stages: [
    { duration: "1m", target: 500 },
    { duration: "1m", target: 1000 },
    // { duration: "1m", target: 3000 },
  ],
  thresholds: {
    http_req_duration: ["p(95)<500"], // 95% of requests should be below 500ms
  },
};

// export default function () {
//   const res = http.get("http://localhost:4000/manga/popular");
//   const res = http.get("http://localhost:4000/manga?name=Vinland%20Saga");
//   // const res = http.get(
//   // //   "https://vivacious-kristyn-chimas-00c6f067.koyeb.app/manga/popular"
//   // // );
//   check(res, { "status was 200": (r) => r.status === 200 });
//   sleep(1);
// }

export default function () {
  const urls = [
    "http://localhost:4000/manga/many",
    "http://localhost:4000/manga/popular",
    "http://localhost:4000/manga?name=Vinland%20Saga",
    "http://localhost:4000/manga?name=Sousou%20no%20Frieren",
  ];

  urls.forEach((url) => {
    const res = http.get(url);
    check(res, { "status was 200": (r) => r.status === 200 });
  });

  sleep(1);
}
