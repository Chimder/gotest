import http from "k6/http";
import { check, sleep } from "k6";
import { Counter } from "k6/metrics";

export const requests = new Counter("http_reqs");

export const options = {
  stages: [
    { duration: "1m", target: 100 }, // Stay at 100 users for 3 minutes
    { duration: "1m", target: 500 }, // Stay at 100 users for 3 minutes
    // { duration: "1m", target: 3000 }, // Stay at 100 users for 3 minutes
  ],
  thresholds: {
    http_req_duration: ["p(95)<500"], // 95% of requests should be below 500ms
  },
};

export default function () {
  const res = http.get("http://localhost:4000/manga?name=Vinland%20Saga");
  // const res = http.get("http://localhost:4000/manga/many");
  check(res, { "status was 200": (r) => r.status === 200 });
  sleep(1);
}

// export default function () {
//   const urls = [
//     "http://localhost:4000/manga/many",
//   ];

//   urls.forEach((url) => {
//     const res = http.get(url);
//     check(res, { "status was 200": (r) => r.status === 200 });
//   });

//   sleep(1);
// }
