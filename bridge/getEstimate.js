const { ethers } = require("ethers");
require('dotenv').config();
const DONconsumerABI = require('../out/DONconsumer.sol/DONconsumer.json');


module.exports = async() => {
    const provider = new ethers.providers.JsonRpcProvider(process.env.MUMBAI_API_KEY);
    const signer = provider.getSigner();

    const number = await provider.getBlockNumber();
    console.log(number);
}

module.exports.tags = ['getEstimate'];
