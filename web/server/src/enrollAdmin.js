'use strict';

const FabricCAServices = require('fabric-ca-client');
const { FileSystemWallet, X509WalletMixin } = require('fabric-network');
const fs = require('fs');
const path = require('path');

// capture network variables from config.json
// const configPath = path.join(process.cwd(), './config.json');
// const configJSON = fs.readFileSync(configPath, 'utf8');
// const config = JSON.parse(configJSON);

// let connection_file = config.connection_file;
// let appAdmin = config.appAdmin;
// let appAdminSecret = config.appAdminSecret;
// let orgMSPID = config.orgMSPID;
// let caName = config.caName;

exports.enrollAdmin = async function (caClient, wallet, orgMspId) {
	try {
			// Check to see if we've already enrolled the admin user.
			const identity = await wallet.get(adminUserId);
			if (identity) {
					console.log('An identity for the admin user already exists in the wallet');
					return;
			}

			// Enroll the admin user, and import the new identity into the wallet.
			const enrollment = await caClient.enroll({ enrollmentID: adminUserId, enrollmentSecret: adminUserPasswd });
			const x509Identity = {
					credentials: {
							certificate: enrollment.certificate,
							privateKey: enrollment.key.toBytes(),
					},
					mspId: orgMspId,
					type: 'X.509',
			};
			await wallet.put(adminUserId, x509Identity);
			console.log('Successfully enrolled admin user and imported it into the wallet');
	} catch (error) {
			console.error(`Failed to enroll admin user : ${error}`);
	}
};
