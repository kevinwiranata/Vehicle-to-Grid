"use strict";
const { Wallets, Gateway } = require("fabric-network");
const path = require("path");
const util = require("util");
const fs = require("fs");
const FabricCAServices = require("fabric-ca-client");
const adminUserId = 'admin';
const enrollAdmin = require('../src/enrollAdmin');

//connect to the config file
// const configPath = path.join(process.cwd(), "./config.json");
// const configJSON = fs.readFileSync(configPath, "utf8");
// const config = JSON.parse(configJSON);
// let connection_file = config.connection_file;
// let gatewayDiscovery = config.gatewayDiscovery;
// let appAdmin = config.appAdmin;
// let orgMSPID = config.orgMSPID;

// connect to the connection file
// const ccpPath = path.join(process.cwd(), connection_file);
// const ccpJSON = fs.readFileSync(ccpPath, "utf8");
// const ccp = JSON.parse(ccpJSON);

// creates a network object (contract, network, gateway) to invoke transactions
exports.invoke = async function (userName, smartContract, isQuery, func, args) {
	const gateway = new Gateway();

	try {
		const walletPath = path.join(process.cwd(), "wallet");
		const wallet = await Wallets.newFileSystemWallet(walletPath);
		console.log(`Wallet path: ${walletPath}`);
		console.log("userName: ");
		console.log(userName);

		console.log("wallet: ");
		console.log(util.inspect(wallet));

		const userExists = await wallet.get(userName);
		if (!userExists) {
			console.log(
				"An identity for the user " + userName + " does not exist in the wallet"
			);
			let response = {};
			response.error =
				"An identity for the user " +
				userName +
				" does not exist in the wallet. Register first";
			return response;
		}

		await gateway.connect(ccp, {
			wallet,
			identity: appAdmin,
			discovery: {enabled: true, asLocalhost:false},
		});

		// Connect to our local fabric
		const network = await gateway.getNetwork("chanel1");
		console.log("Connected to mychannel. ");
		// Get the contract we have installed on the peer
		const contract = await network.getContract(smartContract);

		// invoke here
		if (isQuery) {
			// query, transaction is not recorded on ledger
			if (args) {
				let response = await contract.evaluateTransaction(func, args);
				console.log(response);
				console.log(`Transaction ${func} with args ${args} has been evaluated`);
				await gateway.disconnect();
				return response;
			} else {
				let response = await networkObj.contract.evaluateTransaction(func);
				console.log(response);
				console.log(`Transaction ${func} without args has been evaluated`);
				await gateway.disconnect();
				return response;
			}
		} else {
			if (args) {
				let response = await contract.submitTransaction(func, args);
				console.log(response);
				console.log(`Transaction ${func} with args ${args} has been evaluated`);
				await gateway.disconnect();
				return response;
			} else {
				let response = await networkObj.contract.submitTransaction(func);
				console.log(response);
				console.log(`Transaction ${func} without args has been evaluated`);
				await gateway.disconnect();
				return response;
			}
		}
	} catch (error) {
		console.log(`Error processing transaction. ${error}`);
		console.log(error.stack);
		let response = {};
		response.error = error;
		console.log("Done connecting to network.");
		gateway.disconnect();
		return response;
	}
};

//register EV
exports.registerEV = async function (EVID, model, age) {
	if (!EVID || !model || !age) {
		let response = {};
		response.error =
			"Error! You need to fill all fields before you can register!";
		return response;
	}

	try {
		// Create a new file system based wallet for managing identities.
		const walletPath = path.join(process.cwd(), "wallet");
		const wallet = await Wallets.newFileSystemWallet(walletPath);
		console.log(`Wallet path: ${walletPath}`);
		console.log(wallet);

		// Check to see if we've already enrolled the user.
		const userExists = await wallet.get(EVID);
		if (userExists) {
			let response = {};
			console.log(
				`An identity for the user ${EVID} already exists in the wallet`
			);
			response.error = `Error! An identity for the user ${EVID} already exists in the wallet. Please enter
        a different EVID.`;
			return response;
		}

		// Check to see if we've already enrolled the admin user.
		const adminExists = await wallet.get(appAdmin);
		if (!adminExists) {
			console.log(
				`An identity for the admin user ${appAdmin} does not exist in the wallet`
			);
			console.log(
				'An identity for the admin user "admin" does not exist in the wallet'
			);
			await enrollAdmin(userOrg, ccp);
			adminIdentity = await wallet.get("admin");
			console.log("Admin Enrolled Successfully");
		}

		// Create a new gateway for connecting to our peer node.
		const gateway = new Gateway();
		await gateway.connect(ccp, {
			wallet,
			identity: appAdmin,
			discovery: {enabled: true, asLocalhost:false},
		});

		// Get the CA client object from the gateway for interacting with the CA.
		const ca = gateway.getClient().getCertificateAuthority();
		const adminIdentity = gateway.getCurrentIdentity();
		console.log(`AdminIdentity: + ${adminIdentity}`);

		// build a user object for authenticating with the CA
		const provider = wallet
			.getProviderRegistry()
			.getProvider(adminIdentity.type);
		const adminUser = await provider.getUserContext(adminIdentity, adminUserId);

		// Register the user, enroll the user, and import the new identity into the wallet.
		// if affiliation is specified by client, the affiliation value must be configured in CA
		const secret = await ca.register(
			{
				affiliation,
				enrollmentID: userId,
				role: "client",
			},
			adminUser
		);
		const enrollment = await ca.enroll({
			enrollmentID: EVID,
			enrollmentSecret: secret,
		});
		const x509Identity = {
			credentials: {
				certificate: enrollment.certificate,
				privateKey: enrollment.key.toBytes(),
			},
			mspId: orgMspId,
			type: "X.509",
		};
		await wallet.put(EVID, x509Identity);
		let response = `Successfully registered EV ${EVID}. Use ID ${EVID} to login above.`;
		return response;
	} catch (error) {
		console.error(`Failed to register user + ${EVID} + : ${error}`);
		let response = {};
		response.error = error;
		return response;
	}
};
