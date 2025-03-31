// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract InsightsMarket {
    struct Insight {
        address provider;
        string industry;
        string title;
        string description;
        uint256 price;
        uint256 timestamp;
    }

    mapping(uint256 => Insight) public insights;
    uint256 public insightCount;

    event InsightCreated(
        uint256 indexed id,
        address indexed provider,
        string industry,
        string title,
        uint256 price
    );

    event InsightPurchased(
        uint256 indexed id,
        address indexed buyer,
        address indexed seller,
        uint256 price
    );

    function createInsight(
        string memory _industry,
        string memory _title,
        string memory _description,
        uint256 _price
    ) public {
        require(bytes(_industry).length > 0, "Industry cannot be empty");
        require(bytes(_title).length > 0, "Title cannot be empty");
        require(_price > 0, "Price must be greater than 0");

        uint256 insightId = insightCount++;
        insights[insightId] = Insight({
            provider: msg.sender,
            industry: _industry,
            title: _title,
            description: _description,
            price: _price,
            timestamp: block.timestamp
        });

        emit InsightCreated(insightId, msg.sender, _industry, _title, _price);
    }

    function purchaseInsight(uint256 _insightId) public payable {
        Insight storage insight = insights[_insightId];
        require(msg.value == insight.price, "Incorrect payment amount");
        require(msg.sender != insight.provider, "Cannot purchase own insight");

        payable(insight.provider).transfer(msg.value);
        emit InsightPurchased(_insightId, msg.sender, insight.provider, msg.value);
    }

    function getInsight(uint256 _insightId)
        public
        view
        returns (
            address provider,
            string memory industry,
            string memory title,
            string memory description,
            uint256 price,
            uint256 timestamp
        )
    {
        Insight storage insight = insights[_insightId];
        return (
            insight.provider,
            insight.industry,
            insight.title,
            insight.description,
            insight.price,
            insight.timestamp
        );
    }
}
