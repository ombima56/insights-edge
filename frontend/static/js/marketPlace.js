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
            try {
                const user = await contract.methods.getUser(userAddress).call();
                document.getElementById("username").innerText = user.username || "Guest";
                document.getElementById("user-type").innerText = user.userType || "User";
                document.getElementById("wallet-balance").innerText = `${user.balance} BPT`;
            } catch (error) {
                console.error("Error loading user profile:", error);
            }
        }

        async function createInsight(industry, type, content, price) {
            try {
                await contract.methods.createInsight(industry, type, content, price).send({ from: userAddress });
                alert("Insight created successfully!");
            } catch (error) {
                console.error("Error creating insight:", error);
            }
        }

        async function purchaseInsight(insightId, price) {
            try {
                await contract.methods.purchaseInsight(insightId).send({ from: userAddress, value: price });
                alert("Insight purchased successfully!");
            } catch (error) {
                console.error("Error purchasing insight:", error);
            }
        }

        async function subscribe(plan) {
            try {
                await contract.methods.subscribe(plan).send({ from: userAddress });
                alert(`Subscribed to ${plan} plan!`);
            } catch (error) {
                console.error("Error subscribing:", error);
            }
        }

        document.getElementById("create-insight-form").addEventListener("submit", async function (e) {
            e.preventDefault();
            const industry = document.getElementById("insight-industry").value;
            const type = document.getElementById("insight-type").value;
            const content = document.getElementById("insight-data").value;
            const price = document.getElementById("insight-price").value;
            await createInsight(industry, type, content, price);
        });

        document.querySelectorAll(".purchase-plan").forEach(button => {
            button.addEventListener("click", async function () {
                const plan = this.getAttribute("data-plan");
                await subscribe(plan);
            });
        });

        document.querySelectorAll(".purchase-insight").forEach(button => {
            button.addEventListener("click", async function () {
                const insightId = this.getAttribute("data-id");
                const price = this.getAttribute("data-price");
                await purchaseInsight(insightId, price);
            });
        });

        loadUserProfile();
    } else {
        alert("Please install MetaMask to use this marketplace.");
    }
});
