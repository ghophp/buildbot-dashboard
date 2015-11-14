var BuildWidget = React.createClass({
    propTypes: {
        builder: React.PropTypes.any.isRequired
    },
    getInitialState: function getInitialState() {
        return { status: "build" };
    },
    tick() {
        this.setState({ status: "building" }); console.log('ok');
    },
    render: function() {
        return React.createElement(
            "div",
            { className: "widget new", "data-status": this.state.status, onClick: this.tick.bind(this) },
            React.createElement(
                "h1",
                { className: "title" },
                this.props.builder.id
            ),
            React.createElement(
                "h2",
                { className: "value" },
                this.props.builder.cachedBuilds.length
            ),
            React.createElement(
                "p",
                { className: "more-info" },
                "cached builds"
            ),
            React.createElement(
                "p",
                { className: "updated-at" },
                "Last updated at 12:00"
            ),
            this.props.children
        );
    }
});