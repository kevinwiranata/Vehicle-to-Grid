"use strict";

const express = require("express");
const bodyParser = require("body-parser");
const path = require("path");
const cors = require("cors");
const morgan = require("morgan");

const index = require("./routes/index");
const EVRouter = require("./routes/EVRouter");

const app = express();
app.use(morgan("dev"));
app.use(bodyParser.json());
app.use(cors());

// // view engine setup
// app.set("views", path.join(__dirname, "views"));
// app.set('view engine', 'html');

app.use("/", index);
app.use("/EV", EVRouter);

// catch 404 and forward to error handler
app.use((req, res, next) => {
  var err = new Error("Not Found");
  err.status = 404;
  next(err);
});

// error handler

app.listen(process.env.PORT || 8081);

// //get all assets in world state
// app.get("/queryAll", async (req, res) => {
//   let networkObj = await network.connectToNetwork(appAdmin, "queryContract");
//   let response = await network.invoke(networkObj, true, "querow myAll", "");
//   let parsedResponse = await JSON.parse(response);
//   res.send(parsedResponse);
// });