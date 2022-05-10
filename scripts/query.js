const axios = require("axios");

const sendGetRequest = async () => {
	const res = await axios.get(
		"https://u0ekz2lj7j-u0xeyovsc6-connect.us0-aws-ws.kaleido.io/identities?fly-sync=true",
		{
			auth: {
				username: "u0anvf575c",
				password: "5o7bDtmzd0wpGBJKdppZ8iCYhseUCtxJwoP0NRkQSIA",
			},
		}
	);
	console.log(res);
};

sendGetRequest();
