<script>
	import { Router, Route, Link } from "svelte-routing";
	import Home from "./pages/Home.svelte";
	import Env from "./pages/Env.svelte";
	import EnvPage from "./pages/EnvPage.svelte";

	export let url = "";
</script>

<main>
	<!-- svelte-ignore missing-declaration -->
	<Router {url}>
		<nav>
			<Link to="/">Home</Link>
			<Link to="env">Env</Link>
		</nav>
		<div>
			<Route path="env" component={Env} />
			<!--for now the router just support case sensitive,
			   one workaround colud be add two time the route
			   Example.
			  <Route path="About" component="{About}" /> 
		   -->
			<Route path="/"><Home /></Route>
			<Route path="env/:id" let:params>
				<EnvPage id={params.id} />
			</Route>
		</div>
	</Router>
</main>

<style>
	main {
		text-align: center;
		padding: 1em;
		max-width: 240px;
		margin: 0 auto;
	}

	h1 {
		color: #ff3e00;
		text-transform: uppercase;
		font-size: 4em;
		font-weight: 100;
	}

	@media (min-width: 640px) {
		main {
			max-width: none;
		}
	}
</style>
