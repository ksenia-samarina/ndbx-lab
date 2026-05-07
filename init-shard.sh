#!/bin/bash

sleep 10

echo "--- Инициализация Config Server ---"
mongosh --host "cfg-1:$MONGODB_PORT" --eval '
rs.initiate({
  _id: "cfg-rs",
  configsvr: true,
  members: [{ _id: 0, host: "cfg-1:'"$MONGODB_PORT"'" }]
});
'

echo "--- Инициализация Shard 1 ---"
mongosh --host "shard-1-a:$MONGODB_PORT" --eval '
rs.initiate({
  _id: "shard-1-rs",
  members: [
    { _id: 0, host: "shard-1-a:'"$MONGODB_PORT"'" },
    { _id: 1, host: "shard-1-b:'"$MONGODB_PORT"'" },
    { _id: 2, host: "shard-1-c:'"$MONGODB_PORT"'" }
  ]
});
'

echo "--- Инициализация Shard 2 ---"
mongosh --host "shard-2-a:$MONGODB_PORT" --eval '
rs.initiate({
  _id: "shard-2-rs",
  members: [
    { _id: 0, host: "shard-2-a:'"$MONGODB_PORT"'" },
    { _id: 1, host: "shard-2-b:'"$MONGODB_PORT"'" },
    { _id: 2, host: "shard-2-c:'"$MONGODB_PORT"'" }
  ]
});
'

sleep 5

echo "--- Настройка шардинга ---"
mongosh --host "$MONGODB_HOST:$MONGODB_PORT" --eval "
sh.addShard('shard-1-rs/shard-1-a:$MONGODB_PORT');
sh.addShard('shard-2-rs/shard-2-a:$MONGODB_PORT');

sh.enableSharding('$MONGODB_DATABASE');

sh.shardCollection(
  '$MONGODB_DATABASE.events',
  { created_by: 'hashed' }
);
"