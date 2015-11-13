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
            this.props.children
        );
    }
});