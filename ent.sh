#!/bin/bash

go run entgo.io/ent/cmd/ent generate --feature sql/upsert --feature sql/modifier ./ent/schema
