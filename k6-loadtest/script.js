import http from 'k6/http';
import { check, sleep } from 'k6';
export let options = {
    vus: 10,
    //duration: '60s',
    iterations: 1000
};


export default function () {
    var params = {
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer token'
        },
    };
    var body = JSON.stringify({
        "email": "user@example.com.br",
        "first": "User",
        "last": "Last",
        "status": "active",
        "age": 43
    })

    let res = http.post('http://localhost:8080/v1/user/', body, params)
    //let res = http.post('http://localhost:3000/v1/user/', body, params)
    //console.log(`Status: ${res.status}, Body: ${res.body}`);
    //check(res, { 'status was 200': (r) => r.status == 200 });
    //check(res, { 'status was 401': (r) => r.status == 401 });
    //sleep(1);
}
