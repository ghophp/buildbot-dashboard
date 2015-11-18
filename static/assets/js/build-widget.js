var BuildWidget = React.createClass({
    propTypes: {
        builder: React.PropTypes.any.isRequired
    },
    getInitialState: function getInitialState() {
        return { 
            status: "build",
            last_build: this.getLastBuildNumber(),
            last_update: ''
        };
    },
    getLastBuildNumber() {
        var c = this.props.builder.cachedBuilds;
        return c && c.length ? c[c.length - 1] : 0;
    },
    openDetails() {
        window.open(buildbotUrl + 'builders/' + this.props.builder.id);
    },
    updateBuilder(builder) {
        console.log(builder);
        this.props.builder = builder;

        var now = moment();
        if (this.props.builder.last_update) {
            var t = new String(this.props.builder.last_update);
            now = moment(t.substr(0, t.length - 3), "X");
        }

        this.setState({ 
            status: builder.state, 
            last_build: this.getLastBuildNumber(),
            last_update: now.format('MMMM Do YYYY, HH:mm:ss') 
        });
    },
    render: function() {
        return React.createElement(
            "div",
            { className: "widget new", "data-status": this.state.status, onClick: this.openDetails.bind(this) },
            React.createElement(
                "h1",
                { className: "title" },
                this.props.builder.id
            ),
            React.createElement(
                "h2",
                { className: "value" },
                this.state.last_build
            ),
            React.createElement(
                "p",
                { className: "more-info" },
                "cached builds"
            ),
            React.createElement(
                "p",
                { className: "updated-at" },
                this.state.last_update
            ),
            this.props.children
        );
    }
});