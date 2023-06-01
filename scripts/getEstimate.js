const { ethers } = require("ethers");
const fs = require("fs");
const DONconsumerABI = require('../out/DONconsumer.sol/DONconsumer.json');
require('dotenv').config();

const main = async() => {
  const provider = new ethers.providers.JsonRpcProvider(process.env.MUMBAI_RPC_URL);
  const signer = new ethers.Wallet(process.env.PRIVATE_KEY ?? "", provider);
  const requestConfig = {
    source: fs.readFileSync("./calculation-example.js").toString(),
    args: ["1", "bitcoin", "btc-bitcoin", "btc", "1000000", "450"],
    secrets: { apiKey: process.env.COINMARKETCAP_API_KEY ?? "" },
    subId: process.env.SUB_ID ?? 0,
    gasLimit: process.env.GAS_LIMIT ?? 1000000,
    gasPrice: (await provider.getFeeData()).gasPrice
  }
  const DONconsumerAddress = process.env.DON_CONSUMER_CONTRACT ?? ''; 
  const consumerContract = new ethers.Contract(DONconsumerAddress, DONconsumerABI.abi, signer);
  const simulatedSecretsURLBytes = `0x${Buffer.from(process.env.BUFFER_URL ?? "").toString("hex")}`;

  const estimatedCost = await consumerContract.getCostEstimate(
    requestConfig.source,
    requestConfig.secrets && Object.keys(requestConfig.secrets).length > 0 ? simulatedSecretsURLBytes : [],
    requestConfig.args ?? [],
    requestConfig.subId,
    requestConfig.gasLimit,
    requestConfig.gasPrice
  );
  console.log("estimatedCost: ", estimatedCost.toNumber()/1e8);

  const overrides = {
    gasLimit: 1_500_000
  }

  const executeRes = await consumerContract.executeRequest(
    requestConfig.source,
    requestConfig.secrets && Object.keys(requestConfig.secrets).length > 0 ? simulatedSecretsURLBytes : [],
    requestConfig.args ?? [],
    requestConfig.subId,
    requestConfig.gasLimit,
    overrides
  );

  console.log({executeRes});
} 

main()
