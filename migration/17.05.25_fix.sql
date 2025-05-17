
ALTER TABLE friendships
DROP CONSTRAINT check_username_order,
ADD CONSTRAINT check_username_order CHECK (user1_username < user2_username COLLATE "C");
