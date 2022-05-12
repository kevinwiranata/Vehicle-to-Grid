const axios = require("axios");

const sendGetRequest = async () => {
	await axios
		.post(
			"https://u0ekz2lj7j-u0xeyovsc6-connect.us0-aws-ws.kaleido.io/query?fly-sync=true",
			{
				headers: {
					signer: "user2",
					channel: "default-channel",
					chaincode: "ev_contract",
				},
				func: "QueryAssetHistory",
				args: ["1"],
				strongread: true,
			},
			{
				// firefly basic auth
				auth: {
					username: process.env.USERNAME,
					password: process.env.PASSWORD,
				},
			}
		)
		.then((resp) => {
			console.log(resp.data);
		})
		.catch((err) => {
			// Handle Error Here
			console.error(err);
		});
};
sendGetRequest();
