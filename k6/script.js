import http from 'k6/http';
import { check } from 'k6';

const url = 'http://localhost:8080'

export const options = {
    vus: 50,
    iterations: 10000,
};

export function setup() {
    let res1 = http.post(url + '/accounts', JSON.stringify({ document_number: '1234567' }), {
        headers: { 'Content-Type': 'application/json' },
    });
    let res2 = http.post(url + '/accounts', JSON.stringify({ document_number: '87654321' }), {
        headers: { 'Content-Type': 'application/json' },
    });
    check(res1, {
        'account 1 created': (r) => r.status === 201,
    });
    check(res2, {
        'account 2 created': (r) => r.status === 201,
    });
}

export default function () {
    let total = 0.0
    let amount1 = parseFloat((Math.random() * 1000).toFixed(2));
    let amount2 = parseFloat((Math.random() * 1000).toFixed(2));
    total -= amount1;
    total += amount2;
    let body1 = { account_id: 1, operation_type_id: 1, amount: amount1 };
    let body2 = { account_id: 1, operation_type_id: 4, amount: amount2 };

    let res1 = http.post(url + '/transactions', JSON.stringify(body1), {
        headers: { 'Content-Type': 'application/json' },
    });
    let res2 = http.post(url + '/transactions', JSON.stringify(body2), {
        headers: { 'Content-Type': 'application/json' },
    });

    check(res1, {
        'transaction 1 created': (r) => r.status === 201,
    });
    check(res2, {
        'transaction 2 created': (r) => r.status === 201,
    });

    let optype = 1
    if (total > 0) {
        optype = 4
    }
    total = parseFloat(total.toFixed(2));
    let body3 = { account_id: 2, operation_type_id: optype, amount: total };
    let tx1 = http.post(url + '/transactions', JSON.stringify(body3), {
        headers: { 'Content-Type': 'application/json' },
    });
    check(tx1, {
        'followup transaction created': (r) => r.status === 201,
    });
}

export function teardown() {
    const res = http.get(url + '/accounts/1/balance');
    check(res, {
        'is status 200': (r) => r.status === 200,
    });
    const res2 = http.get(url + '/accounts/2/balance');
    check(res2, {
        'is status 200': (r) => r.status === 200,
    });
    check({ b1: res.json().balance, b2: res2.json().balance }, {
        'balance is correct': (r) => r.b1 === r.b2,
    })
}