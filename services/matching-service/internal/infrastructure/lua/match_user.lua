-- ARGV[1]: current user ID
-- ARGV[2]: max scan limit
-- KEYS: list of tag bitmap keys

local userID = ARGV[1]
local maxScan = tonumber(ARGV[2])

-- Get or assign compact index
local currentIndex = redis.call("HGET", "user:index", userID)
if not currentIndex then
    currentIndex = redis.call("INCR", "user:index:counter")
    redis.call("HSET", "user:index", userID, currentIndex)
    redis.call("HSET", "index:user", currentIndex, userID)
end
currentIndex = tonumber(currentIndex)

-- Return -1 if no tag keys provided
if #KEYS == 0 then
    return -1
end

-- Temporary key for AND result
local andKey = "temp:and:" .. currentIndex .. ":" .. redis.call("TIME")[1]
redis.call("DEL", andKey)

-- Perform BITOP AND on all tag keys
redis.call("BITOP", "AND", andKey, unpack(KEYS))

-- Find the first set bit
local matchedIndex = redis.call("BITPOS", andKey, 1)

-- Clean up temp key
redis.call("DEL", andKey)

-- No match found or out of bounds
if matchedIndex == -1 or matchedIndex >= maxScan then
    return -1
end

-- Avoid matching with self
if matchedIndex == currentIndex then
    return -1
end

-- Remove both users from tag bitmaps
for _, key in ipairs(KEYS) do
    redis.call("SETBIT", key, currentIndex, 0)
    redis.call("SETBIT", key, matchedIndex, 0)
end

-- Return matched user ID (not internal index)
local matchedUserID = redis.call("HGET", "index:user", matchedIndex)
if not matchedUserID then
    return -1
end

return matchedUserID
