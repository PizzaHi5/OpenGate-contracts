//const { ethers } = require("ethers");
const hre = require("hardhat");

//const { ethers } = require("hardhat");
require('dotenv').config();
const DONconsumerABI = require('../out/DONconsumer.sol/DONconsumer.json');


module.exports = async() => {
    const provider = new hre.ethers.providers.JsonRpcProvider(process.env.MUMBAI_API_KEY);
    const signer = hre.provider.getSigner();

    const number = await hre.provider.getBlockNumber();
    console.log(number);
}

module.exports.tags = ['getEstimate'];
