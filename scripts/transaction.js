const axios = require("axios");

const data = {
	EV1: {EVID: "ID1", CSOID: "CSO1", PowerFlow: "24.8", RecentMoney: "2.43", Temperature: 45.2, SoC: 87.2,},
	EV2: {EVID: "ID2", CSOID: "CSO1", PowerFlow: "24.8",RecentMoney: "2.43",Temperature: 45.2,SoC: 87.2,},
	EV3: {EVID: "ID3",CSOID: "CSO1",PowerFlow: "24.8",RecentMoney: "2.43",Temperature: 45.2,SoC: 87.2},
	EV4: {EVID: "ID4",CSOID: "CSO1",PowerFlow: "24.8",RecentMoney: "2.43",Temperature: 45.2,SoC: 87.2},
	EV5: {EVID: "ID5",CSOID: "CSO1",PowerFlow: "24.8",RecentMoney: "2.43",Temperature: 45.2,SoC: 87.2},
	EV6: {EVID: "ID6",CSOID: "CSO1",PowerFlow: "24.8",RecentMoney: "2.43",Temperature: 45.2,SoC: 87.2},
};



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
						"CSO1",
						String(data[EV].PowerFlow),
						String(data[EV].RecentMoney),
						String(data[EV].Temperature),
						String(data[EV].SoC),
					],
					init: false,
				}, {
					// firefly basic auth 
					auth: {
						username: process.env.username,
						password: process.env.password
					}
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

sendTx(data);
