#!/usr/bin/env python

import sqlite3
import sys
import os

def initdb():
    db_path = "/tmp/test.db"
    conn = None

    if not os.path.exists(db_path):
        conn = sqlite3.connect(db_path)
        c = conn.cursor()
        c.execute("""CREATE TABLE IF NOT EXISTS `tweets` (
`TweetID` INTEGER,
`UserID` INTEGER,
`tweet` TEXT,
`date` TEXT,
PRIMARY KEY(`TweetID`)
);""")

        c.execute("""CREATE TABLE IF NOT EXISTS `users` (
     `UserID` INTEGER,
     `login` TEXT,
     `password` TEXT,
     `salt` TEXT,
     PRIMARY KEY(`UserID`)
    );
    """)
        conn.commit()
        conn.close()

    conn = sqlite3.connect(db_path)
    c = conn.cursor()
    return (conn,c)




#c = conn.execute("select MessageID, Subject FROM messages WHERE MessageID = '{}'".format(sys.argv[1]))
#
#
#conn.close()


def main(args):
    conn, c = initdb()
    if(args.register):
        username = args.register.split(":")[0] #args.register[:args.register.index(":")]
        password = args.register.split(":")[1] #args.register[1+args.register.index(":"):]
        salt = args.register.split(":")[2] #args.register[1+args.register.index(":"):]
        c.execute("insert into `users`('login', 'password', 'salt') VALUES('{}','{}','{}');".format(username, password, salt))
        conn.commit()
    elif(args.get_user_id):
        c2 = c.execute("select UserID, login from `users` WHERE login ='{}';".format(args.get_user_id))
        res = c2.fetchall()
        if(len(res ) > 0):
            print "{}".format(res[0][0])
    elif(args.get_uuid_password):
        c2 = c.execute("select salt from `users` WHERE UserID ='{}';".format(args.get_uuid_password))
        res = c2.fetchall()
        if(len(res ) > 0):
            print "{}".format(res[0][0])
    elif(args.user_already_exist):
        username = args.user_already_exist.split(":")[0]
        password = args.user_already_exist.split(":")[1]
        c2 = c.execute("select login, password from `users` WHERE login ='{}' AND password ='{}' ;".format(username, password))
        res = c2.fetchall()
        if(len(res ) > 0):
            print "{}".format(res[0][0])
    elif(args.get_user_password):
        c2 = c.execute("select password,UserID from `users` WHERE UserID ='{}';".format(args.get_user_password))
        res = c2.fetchall()
        if(len(res) > 0):
            print "{}".format(res[0][0])
    elif(args.get_user_name):
        c2 = c.execute("select login,UserID from `users` WHERE UserID ='{}';".format(args.get_user_name))
        res = c2.fetchall()
        if(len(res) > 0):
            print "{}".format(res[0][0])
    elif(args.register_tweet):
        userid = args.register_tweet[:args.register_tweet.index(":")]
        tweet_content = args.register_tweet[1+args.register_tweet.index(":"):]
        c.execute("insert into `tweets`('UserID', 'tweet', 'date') VALUES('{}','{}', datetime('now'));".format(userid, tweet_content ))
        conn.commit()
    elif(args.get_tweet_content):
        c2 = c.execute("select tweet,TweetId from `tweets` WHERE TweetId ='{}';".format(args.get_tweet_content))
        res = c2.fetchall()
        if(len(res) > 0):
            print "{}".format(res[0][0])
    elif(args.get_tweet_date):
        c2 = c.execute("select date,TweetId from `tweets` WHERE TweetId ='{}';".format(args.get_tweet_date))
        res = c2.fetchall()
        if(len(res) > 0):
            print "{}".format(res[0][0])
    elif(args.get_tweet_userid):
        c2 = c.execute("select UserID,TweetId from `tweets` WHERE TweetId ='{}';".format(args.get_tweet_userid))
        res = c2.fetchall()
        if(len(res) > 0):
            print "{}".format(res[0][0])
    elif(args.list_tweet_ids):
        c2 = c.execute("select TweetId, UserID from `tweets` WHERE UserID ='{}';".format(args.list_tweet_ids))
        res = c2.fetchall()
        if(len(res) > 0):
            print ",".join(map(lambda v:"{}".format(v[0]),res))
    elif(args.lastid):
        c2 = c.execute("SELECT rowid from {} order by ROWID DESC limit 1;".format(args.lastid))
        res = c2.fetchall()
        if(len(res ) > 0):
            print "{}".format(res[0][0])
    conn.close()


if __name__=="__main__":
    import argparse
    parser = argparse.ArgumentParser()
    group = parser.add_mutually_exclusive_group()
    group.add_argument("--get_user_id", help="return user_id for a given username, return nothing if there is no user with this name", type=str)
    group.add_argument("--get_user_password", help="return password for a given user_id, return nothing if there is no user with this user_id", type=str)
    group.add_argument("--get_user_name", help="return name for a given user_id, return nothing if there is no user with this user_id", type=str)
    group.add_argument("--register", help="register a new user. You must specify a string with 'username:password' the first ':' character will be the delimiter", type=str)

    group.add_argument("--lastid", help="get the last id for the users's table",type=str)
    group.add_argument("--get_uuid_password", help="get the uuid for the user",type=str)
    group.add_argument("--user_already_exist", help="verify that user isn't present",type=str)

    group.add_argument("--get_tweet_content", help="get tweet content for a given tweet id, return nothing if nothing found for this id", type=str)
    group.add_argument("--get_tweet_date", help="get tweet date for a given tweet id, return nothing if nothing found for this id", type=str)
    group.add_argument("--get_tweet_userid", help="get tweet userid for a given tweet id, return nothing if nothing found for this id", type=str)
    group.add_argument("--list_tweet_ids", help="get all tweet ids for a given userid, tweetsid will be comma separated", type=str)

    group.add_argument("--register_tweet", help="register new tweet, given a string in the form of 'userid:tweet_content' ", type=str)

    main(parser.parse_args())
