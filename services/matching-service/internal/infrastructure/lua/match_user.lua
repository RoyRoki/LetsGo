-- ARGV[1]: current user ID
-- ARGV[2]: max scan limit

local currentUserID = tonumber(ARGV[1])
local maxScan = tonumber(ARGV[2])

-- Temporary key for AND result
local andKey = "temp:and:" .. currentUserID .. ":" .. redis.call("TIME")[1]
redis.call("DEL", andKey)

-- Perform BITOP AND on all keys
redis.call("BITOP", "AND", andKey, unpack(KEYS))

-- Find first set bit
local matchedIndex = redis.call("BITPOS", andKey, 1)

-- Clean up temp key
redis.call("DEL", andKey)

-- No match found or out of bounds
if matchedIndex == -1 or matchedIndex >= maxScan then
	return -1
end

-- Don't match with self
if matchedIndex == currentUserID then
	return -1
end

-- Remove both users from all tag bitmaps
for _, key in ipairs(KEYS) do
	redis.call("SETBIT", key, currentUserID, 0)
	redis.call("SETBIT", key, matchedIndex, 0)
end

return matchedIndex
