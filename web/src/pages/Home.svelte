<script>
    import { postData } from "../util.js";

    let promise = postData("http://localhost:9999/collections/all", {
        skip: 0,
        limit: 10,
    }).then((data) => {
        return data; // JSON data parsed by `data.json()` call
    });
</script>

{#await promise}
    <!-- optionally show something while promise is pending -->
    <p>Loading..</p>
{:then data}
    <!-- promise was fulfilled -->
    {#each data.body as colls}
        <div class="container">
            <div class="env-list">
                <div class="coll">
                    {colls.name}
                </div>
                <div>
                    Description: {colls.info.description}
                </div>
                <div>
                    Resources: {colls.item[0].item.length}
                </div>
                <div>Descendants:</div>
                <div>
                    {#each colls.descendants as desc}
                        <div>
                            Name: {desc.Name}
                        </div>

                        <div>
                            <ul>
                                {#each desc.Resources as res}
                                    <li>{res}</li>
                                {/each}
                            </ul>
                        </div>
                    {/each}
                </div>
            </div>
        </div>
    {/each}
{:catch error}
    <!-- optionally show something while promise was rejected -->
{/await}

<style>
    .container {
        display: flex;
        justify-items: start;
        justify-content: start;
        border-bottom: 0.5px solid black;
        max-height: 300px;
        overflow-x: scroll;
    }

    .env-list {
        padding: 10px;
        margin: 20px;
        display: grid;
        justify-items: start;
        justify-content: start;
        /* box-shadow: 10px 10px 5px #888888; */
    }

    .coll {
        font-weight: 600;
    }

    .endpoint {
        padding: 5px;
    }

    .env-list:hover {
        cursor: pointer;
    }
</style>
