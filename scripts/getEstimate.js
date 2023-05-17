//const { ethers } = require("ethers");
//const hre = require("@nomicfoundation/hardhat-foundry");


const hree = require("hardhat");
require('dotenv').config();
const DONconsumerABI = require('../out/DONconsumer.sol/DONconsumer.json');


module.exports = async() => {
    const provider = new hree.ethers.providers.JsonRpcProvider(process.env.MUMBAI_API_KEY);
    const signer = hree.provider.getSigner();

    const number = await hree.provider.getBlockNumber(); 
    console.log(number);

    //const DONaddress = await 

}

module.exports.tags = ['getEstimate'];
