import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
    stages: [
        { duration: '5s', target: 5 }, // Ramp-up to 10 virtual users over 1 minute
        { duration: '10s', target: 10 }, // Stay at 10 virtual users for 3 minutes
        { duration: '15s', target: 20 }, // Ramp-up to 50 virtual users over 1 minute
        { duration: '10s', target: 15 }, // Stay at 50 virtual users for 5 minutes
    ],
    thresholds: {
        'http_req_duration': ['p(95)<500'], // 95% of requests should complete within 500ms
    },
};

function generateRandomIpAddress() {
    // Generate a random number between 0 and 255 for each octet
    const octet1 = 169//Math.floor(Math.random() * 256);
    const octet2 = 1//Math.floor(Math.random() * 256);
    const octet3 = 245//Math.floor(Math.random() * 256);
    const octet4 = Math.floor(Math.random() * 256);

    // Return the formatted IP address
    return `${octet1}.${octet2}.${octet3}.${octet4}`;
}

export default function () {
    // Define an array to store the generated IP addresses
    const ipAddresses = [];

    // Generate a collection of IP addresses
    for (let i = 0; i < 100; i++) {
        // Generate a random IP address
        const ipAddress = generateRandomIpAddress();

        // Add the IP address to the array
        ipAddresses.push(ipAddress);
    }

    // Loop through the IP addresses and make requests
    for (const ipAddress of ipAddresses) {
        // Make a GET request to the URL with the IP address
        const res = http.get(`http://127.0.0.1:3000/api/lookup/${ipAddress}`);

        // Check if the response is successful
        check(res, {
            'is status 200': (r) => r.status === 200,
        });

        // Sleep for a short duration between requests
        sleep(0.1);
    }
}