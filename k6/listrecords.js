import http from 'k6/http';

export default () => {
	let payload = {
    "from": "2021-12-1T01:59:36.764428Z",
    "page_size": 100,
    "to": "2021-12-15T02:12:17.243489800Z"
	}
  http.post('http://localhost:50101/v1/dogfood/records', JSON.stringify(payload), {'Content-Type': 'application/json'});
};
