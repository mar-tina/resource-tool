<script>
    import { onMount } from "svelte";
    import { postData } from "../util";

    export let id;
    let promise = postData("http://localhost:9999/env/fetch", { id: id })
        .then((data) => {
            console.log("data", data);
            return data;
        })
        .catch((e) => {
            console.log("err", e);
        });

    function update(promise) {
        console.log("update", promise);
        postData("http://localhost:9999/env/update", {
            environment: promise.body.environment,
        }).then((data) => {
            console.log("dataaaa", data);
        });
    }
</script>

{#await promise}
    <!-- optionally show something while promise is pending -->
    <p>Loading..</p>
{:then data}
    <!-- promise was fulfilled -->
    <p>{data.body.environment.name}</p>
    {#each data.body.environment.values as val}
        <div class="env-list">
            <input bind:value={val.key} />
            <input bind:value={val.value} />
            <button on:click={() => update(data)}> Update </button>
            <br />
        </div>
    {/each}
{:catch error}
    <!-- optionally show something while promise was rejected -->
{/await}
