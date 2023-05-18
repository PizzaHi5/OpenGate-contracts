const { ethers } = require("ethers");
const fs = require("fs");
const DONconsumerABI = require('../out/DONconsumer.sol/DONconsumer.json');
require('dotenv').config();

const main = async() => {
  const provider = new ethers.providers.JsonRpcProvider(process.env.MUMBAI_RPC_URL);
  const signer = provider.getSigner();
  const requestConfig = {
    source: fs.readFileSync("./calculation-example.js").toString(),
    args: ["1", "bitcoin", "btc-bitcoin", "btc", "1000000", "450"],
    secrets: { apiKey: process.env.COINMARKETCAP_API_KEY ?? "" },
    subId: 655,
    gasLimit: 1000000,
    gasPrice: await (await provider.getFeeData()).gasPrice
  }
  const DONconsumerAddress = process.env.DON_CONSUMER_CONTRACT ?? ''; 
  console.log({DONconsumerAddress})
  const consumerContract = new ethers.Contract(DONconsumerAddress, DONconsumerABI.abi, provider);
  const simulatedSecretsURLBytes = `0x${Buffer.from(
    "https://exampleSecretsURL.com/f09fa0db8d1c8fab8861ec97b1d7fdf1/raw/d49bbd20dc562f035bdf8832399886baefa970c9/encrypted-functions-request-data-1679941580875.json"
  ).toString("hex")}`;

  const estimatedCost = await consumerContract.getCostEstimate(
    requestConfig.source,
    requestConfig.secrets && Object.keys(requestConfig.secrets).length > 0 ? simulatedSecretsURLBytes : [],
    requestConfig.args ?? [],
    requestConfig.subId,
    requestConfig.gasLimit,
    requestConfig.gasPrice
  );

  console.log({estimatedCost});
  console.log(ethers.utils.formatEther(estimatedCost));
} 

main()
