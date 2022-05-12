require('dotenv').config()
const axios = require("axios");
const fs = require("fs");
const path = require("path");
const csv = require("fast-csv");

data = [];
fs.createReadStream(path.resolve(__dirname, "../data", "V2G_one_week.csv"))
	.pipe(csv.parse({ headers: true }))
	.on("error", (error) => console.error(error))
	.on("data", (row) => data.push(row))
	.on("end", () => sendTx(data));

const sendTx = async (data) => {
	for (const EV in data) {
		await axios
			.post(
				"https://u0ekz2lj7j-u0xeyovsc6-connect.us0-aws-ws.kaleido.io/transactions?fly-sync=false",
				{
					headers: {
						type: "SendTransaction",
						signer: "user2",
						channel: "default-channel",
						chaincode: "ev_contract",
					},
					func: "UpdateEVData",
					args: [
						String(data[EV].EVID),
						String(data[EV].CSOID),
						String(data[EV].ChID),
						String(data[EV].Power_flow),
						String(data[EV].Money),
						String(data[EV].Temp),
						String(data[EV].SoC),
						String(data[EV].SoH),
					],
					init: false,
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
	}
};