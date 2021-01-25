<script>
    import { onMount } from "svelte";
    import { postData } from "../util.js";
    import { Link } from "svelte-routing";

    let promise = postData("http://localhost:9999/env/all").then((data) => {
        return data; // JSON data parsed by `data.json()` call
    });
</script>

{#await promise}
    <!-- optionally show something while promise is pending -->
    <p>Loading..</p>
{:then data}
    <!-- promise was fulfilled -->
    {#each data.body.environments as env}
        <div class="env-list">
            <Link to={`env/${env.id}`}>
                <p>{env.name}</p>
            </Link>
        </div>
    {/each}
{:catch error}
    <!-- optionally show something while promise was rejected -->
{/await}

<style>
    .env-list {
        border: 1px solid black;
        padding: 10px;
        margin: 20px;
    }

    .env-list:hover {
        cursor: pointer;
    }
</style>
