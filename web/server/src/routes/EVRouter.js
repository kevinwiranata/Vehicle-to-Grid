"use strict";
const express = require("express");
const path = require("path");
const fs = require("fs");
const bodyParser = require("body-parser");
const network = require("../network.js");

const EVRouter = express.Router();

// const configPath = path.join(process.cwd(), "./config.json");
// const configJSON = fs.readFileSync(configPath, "utf8");
// const config = JSON.parse(configJSON);

// //use this identity to query
// const appAdmin = config.appAdmin;

EVRouter.use(bodyParser.json());

//get EV info, create docot object, and update state with their EVID
EVRouter.route("/register").post(async (req, res) => {
  let EVID = req.body.ID;
  let model = req.body.model;
  let age = req.body.age;

  //first create the identity for the EV and add to wallet
  let response = await network.registerEV(EVID, model, age);
  console.log("response from registerEV: ");
  console.log(response);
  if (response.error) {
    res.statusCode = 500;
    res.send(response.error);
    return;
  } 
  let args = [JSON.stringify(req.body)];

    //connect to network and update the state with EVID
  let invokeResponse = await network.invoke(EVID, "EVContract", false, "createEV", args);

  if (invokeResponse.error) {
    res.statusCode = 500;
    res.send(invokeResponse.error);
    return;
  } else {
    let parsedResponse = JSON.parse(invokeResponse);
    parsedResponse += ". Use EVID to login above.";
    res.send(parsedResponse);
  }
  }
);

EVRouter.route("/login").post(async (req, res) => {
  let networkObj = await network.connectToNetwork(
    req.body.EVID,
    "EVContract"
  );
  if (networkObj.error) {
    res.statusCode = 500;
    res.send(networkObj.error);
    return;
  }

  let args = [JSON.stringify(req.body)];
  //connect to network and update the state with EVID
  let invokeResponse = await network.invoke(
    networkObj,
    true,
    "EVUserExists",
    args
  );

  if (invokeResponse.error) {
    res.statusCode = 500;
    res.send(invokeResponse.error);
    return;
  } else {
    let parsedResponse = JSON.parse(invokeResponse);
    res.send(parsedResponse);
  }
});

module.exports = EVRouter;
