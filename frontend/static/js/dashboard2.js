document.addEventListener("DOMContentLoaded", async function () {
    if (typeof window.ethereum !== "undefined") {
        window.web3 = new Web3(window.ethereum);
        await window.ethereum.request({ method: "eth_requestAccounts" });

        const contractAddress = "YOUR_CONTRACT_ADDRESS";
        const abi = [ /* Your contract ABI here */ ];

        const contract = new web3.eth.Contract(abi, contractAddress);
        const accounts = await web3.eth.getAccounts();
        const userAddress = accounts[0];

        async function loadUserProfile() {
            const user = await contract.methods.getUser(userAddress).call();
            document.getElementById("username").innerText = user.username || "Guest";
            document.getElementById("user-type").innerText = user.userType || "User";
            document.getElementById("wallet-balance").innerText = `${user.balance} BPT`;
        }

        async function createInsight(industry, type, content) {
            await contract.methods.createInsight(industry, type, content).send({ from: userAddress });
            alert("Insight created successfully!");
        }

        async function subscribe(plan) {
            await contract.methods.subscribe(plan).send({ from: userAddress });
            alert(`Subscribed to ${plan} plan!`);
        }

        document.getElementById("create-insight-form").addEventListener("submit", async function (e) {
            e.preventDefault();
            const industry = document.getElementById("insight-industry").value;
            const type = document.getElementById("insight-type").value;
            const content = document.getElementById("insight-data").value;
            await createInsight(industry, type, content);
        });

        document.querySelectorAll(".purchase-plan").forEach(button => {
            button.addEventListener("click", async function () {
                const plan = this.getAttribute("data-plan");
                await subscribe(plan);
            });
        });

        loadUserProfile();
    } else {
        alert("Please install MetaMask to use this dashboard.");
    }
});
