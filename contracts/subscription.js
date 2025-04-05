document.addEventListener("DOMContentLoaded", async () => {
    if (typeof window.ethereum !== "undefined") {
      const web3 = new Web3(window.ethereum);
      await window.ethereum.request({ method: "eth_requestAccounts" });
  
      const contractAddress = "YOUR_CONTRACT_ADDRESS"; // replace with deployed address
      const abi = [ /* Paste ABI from SubscriptionMarket here */ ];
  
      const contract = new web3.eth.Contract(abi, contractAddress);
      const accounts = await web3.eth.getAccounts();
      const user = accounts[0];
  
      const plans = {
        Monthly: 0,
        Quarterly: 1,
        Yearly: 2
      };
  
      const planPrices = {
        0: web3.utils.toWei("0.01", "ether"),
        1: web3.utils.toWei("0.025", "ether"),
        2: web3.utils.toWei("0.08", "ether")
      };
  
      async function getSubscription() {
        try {
          const result = await contract.methods.getSubscription(user).call();
          const planName = Object.keys(plans)[result.plan] || "None";
          const expiryDate = new Date(result.expiry * 1000).toLocaleString();
          const active = result.active ? "Active" : "Expired";
  
          document.getElementById("sub-plan").innerText = planName;
          document.getElementById("sub-expiry").innerText = expiryDate;
          document.getElementById("sub-status").innerText = active;
        } catch (error) {
          console.error("Error fetching subscription:", error);
        }
      }
  
      async function subscribe(planKey) {
        const planId = plans[planKey];
        const price = planPrices[planId];
  
        try {
          await contract.methods.subscribe(planId).send({
            from: user,
            value: price
          });
          alert(`Subscribed to ${planKey} plan!`);
          getSubscription();
        } catch (error) {
          console.error("Subscription failed:", error);
          alert("Subscription failed.");
        }
      }
  
      // Example buttons:
      document.getElementById("subscribe-monthly").addEventListener("click", () => {
        subscribe("Monthly");
      });
  
      document.getElementById("subscribe-quarterly").addEventListener("click", () => {
        subscribe("Quarterly");
      });
  
      document.getElementById("subscribe-yearly").addEventListener("click", () => {
        subscribe("Yearly");
      });
  
      getSubscription();
    } else {
      alert("Please install MetaMask to use this feature.");
    }
  });