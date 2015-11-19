var LoadingWidget = React.createClass({
    propTypes: {
        
    },
    componentDidMount: function() {
		var el = ReactDOM.findDOMNode(this);

		var $red = el.querySelector('.red');
		var $green = el.querySelector('.green');
		var $blue = el.querySelector('.blue');
		var $purple = el.querySelector('.purple');
		
		var start = 0;

		var tl = new TimelineMax({
			repeat: -1,
			repeatDelay: 0.05,
		});

		var movement = 30;
		if (genericSize === "small") {
			movement = 15;
		}

		tl.fromTo($green, 0.4, {x: 0}, {x: movement, ease: Linear.easeNone}, start);
		tl.fromTo($blue, 0.4, {x: 0}, {x: movement * -1, ease: Linear.easeNone}, start);
		
		tl.to($green, 0.4, {x: movement * 2, ease: Linear.easeNone}, start + 0.45);
		tl.fromTo($purple, 0.4, {x: 0}, {x: movement  * -1, ease: Linear.easeNone}, start + 0.45);
		tl.to($blue, 0.4, {x: (movement * 2) * -1, ease: Linear.easeNone}, start + 0.45);
		tl.fromTo($red, 0.4, {x: 0}, {x: movement, ease: Linear.easeNone}, start + 0.45);
		
		tl.to($green, 0.4, {x: movement, ease: Linear.easeNone}, start + 0.9);
		tl.to($purple, 0.4, {x: 0, ease: Linear.easeNone}, start + 0.9);
		tl.to($blue, 0.4, {x: movement  * -1, ease: Linear.easeNone}, start + 0.9);
		tl.to($red, 0.4, {x: 0, ease: Linear.easeNone}, start + 0.9);
		
		tl.to($green, 0.4, {x: 0, ease: Linear.easeNone}, start + 1.35);
		tl.to($blue, 0.4, {x: 0, ease: Linear.easeNone}, start + 1.35);
	},
    render: function() {
        return React.createElement(
            "div",
            { className: "loading-container" },
            React.createElement(
                "div",
                { className: "dot red" }
            ),
            React.createElement(
                "div",
                { className: "dot green" }
            ),
            React.createElement(
                "div",
                { className: "dot blue" }
            ),
            React.createElement(
                "div",
                { className: "dot purple" }
            ),
            this.props.children
        );
    }
});