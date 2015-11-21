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
    getConfigUrl() {
        return buildbotUrl + "builders/"+ this.props.builder.id +"/";
    },
    getLastBuildUrl() {
        return this.getConfigUrl() +"builds/" + this.getLastBuildNumber();
    },
    render: function() {

        var loadingEl = null;
        if (this.state.status === 'building') {
            loadingEl = React.createElement(LoadingWidget, {});
        }

        var lastBuildEl = null;
        if (this.state.last_build > 0) {
            lastBuildEl = React.createElement(
                "a",
                { className: "lnr lnr-history", href: this.getLastBuildUrl(), target: "_blank" },
                ""
            );
        }

        return React.createElement(
            "div",
            { className: "widget new", "data-status": this.state.status },
            React.createElement(
                "div",
                { className: "icons-wrapper" },
                React.createElement(
                    "a",
                    { className: "lnr lnr-cog", href: this.getConfigUrl(), target: "_blank" },
                    ""
                ),
                lastBuildEl
            ),
            React.createElement(
                "h1",
                { className: "title" },
                this.props.builder.id
            ),
            React.createElement(
                "h2",
                { className: "value" },
                this.state.last_build > 0 ? this.state.last_build : '-'
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
            loadingEl,
            this.props.children
        );
    }
});