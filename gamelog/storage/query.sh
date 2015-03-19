#!/bin/sh
data=`date +%Y%m%d`

dir=./all_$data

diruser=./user_cangku_fz_$data
diruser1=./user_cangku_sx_$data
diruser2=./user_cangku_zhenbao_$data
diruser3=./user_cangku_fujia_$data
fz_config=./fz_setup.ini
sx_config=./sx_setup.ini
zb_config=./zb_setup.ini
fj_config=./fj_setup.ini

username="sgcuser"
password="gamlaxy"
db="testStorage"
table="storage"



rm -f $dir

#SELECT objectId,sum(count) FROM storage group by objectId;

ret=`mysql -u$username -p$password $db -e "SELECT objectId,count(id) FROM storage group by objectId"`

echo $ret|awk '{ i = 3; while ( i <= NF ) { print $i"\t"$(i+1); i+=2}}' >> $dir 


#cat $fz_config|while read line
#do        
#        echo -ne $line"\t" >>$diruser 
#done
#echo -ne "\n" >>$diruser
#cat $fz_config|while read line
#do
#	ret=`mysql -u$username -p$password $db -e "SELECT count(id) FROM storage where objectid=$line"`
#	cnt=`echo $ret|awk '{print $2}'`
#	sum=`mysql -u$username -p$password $db -e "SELECT sum(count) FROM storage where objectid=$line"`
#	echo -ne $cnt"\t" >>$diruser 
#done

