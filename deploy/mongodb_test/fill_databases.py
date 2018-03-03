#!/usr/bin/python2
from pymongo import MongoClient
from multiprocessing import Pool
import random, string

connection = MongoClient('localhost', maxPoolSize=200)

def insert(numbase):
        print(numbase)
	db_name = 'backup_me_{}'.format(numbase)
	db = connection[db_name]
	bulk = list()
	for num_coll in xrange(0,2):
		coll_name = 'user_data_{}'.format(num_coll)
		db[coll_name].insert({"1":''.join(random.choice(string.ascii_uppercase + string.digits) for _ in xrange(20))})



if __name__ == '__main__':
	p = Pool(12)
	p.map(insert, [x for x in xrange(0,300)])
	p.close()
	p.join()

