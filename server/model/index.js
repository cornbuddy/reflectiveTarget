async function initDatabase(client, databaseName) {
    const res = await client.query(
        `SELECT FROM pg_database WHERE datname = '${databaseName}'`,
    );

    if (res.rowCount === 0) {
        console.log(`${databaseName} database not found, creating it.`);
        await client.query(`CREATE DATABASE "${databaseName}";`);
        console.log(`created database ${databaseName}.`);
    } else {
        console.log(`${databaseName} database already exists.`);
    }
}

module.exports = { initDatabase };
