**Procedure for DB**

delimiter
//
-- Create the Procedure called 'GenerateTableAndGetRecord()'"
create procedure GenerateTableAndGetRecord(
    in tableName varchar(255),
    in filePath longtext
)
begin
    -- Content for Procedure 'GenerateTableAndGetRecord()' in here
    set@fields= 'CODE            varchar(42) not null,
        TIME           varchar(5)  not null,
        SAVE_DATE          date        not null,
        SAVE_TIME          varchar(15) not null,
        SERIALNUMBER int         not null,
        FIRST         varchar(2)  not null,
        LAST           varchar(5)  not null,
        DATE            varchar(10) null,
        COLUMN_1          varchar(24) null,
        COLUMN_2          varchar(4)  null,
        COLUMN_3          varchar(4)  null,
        COLUMN_4          varchar(5)  null,
        COLUMN_5          varchar(24) null,
        COLUMN_6          varchar(5)  null,
        COLUMN_7          varchar(24) null,
        COLUMN_8          varchar(5)  null,
        COLUMN_9          varchar(24) null,
        COLUMN_10          varchar(5)  null,
        COLUMN_11          varchar(24) null,
        COLUMN_12          varchar(4)  null,
        COLUMN_13          varchar(24) null,
        COLUMN_14          varchar(5)  null,
        COLUMN_15          varchar(24) null,
        COLUMN_16          varchar(5)  null,
        COLUMN_17          varchar(24) null,
        COLUMN_18          varchar(5)  null,
        COLUMN_19          varchar(4)  null,
        COLUMN_20          varchar(24) null,
        COLUMN_21          varchar(5)  null,
        COLUMN_22          varchar(24) null,
        COLUMN_23          varchar(5)  null,
        COLUMN_24          varchar(24) null,
        COLUMN_25          varchar(5)  null,
        COLUMN_26         varchar(5)  null,
        COLUMN_27          varchar(24) null,
        COLUMN_28          varchar(5)  null,
        COLUMN_29          varchar(5)  null,
        COLUMN_30          varchar(4)  null,
        COLUMN_31          varchar(4)  null,
        COLUMN_32          varchar(24) null,
        COLUMN_33         varchar(5)  null,
        COLUMN_34          varchar(24) null,
        COLUMN_35          varchar(8)  null,
        COLUMN_36          varchar(8)  null,
        COLUMN_37          varchar(8)  null,
        COLUMN_38          varchar(8)  null,
        COLUMN_39          varchar(10) null,
        COLUMN_40          varchar(24) null,
        COLUMN_41         varchar(5)  null,
        COLUMN_42          varchar(24) null,
        COLUMN_43         varchar(5)  null,
        COLUMN_44          varchar(24) null,
        COLUMN_45          varchar(24) null,
        COLUMN_46         varchar(5)  null,
        COLUMN_47          varchar(24) null,
        COLUMN_48          varchar(5)  null,
        COLUMN_49          varchar(24) null,
        COLUMN_50          varchar(5)  null,
        COLUMN_51          varchar(24) null,
        COLUMN_52          varchar(5)  null,
        COLUMN_53          varchar(24) null,
        COLUMN_54          varchar(5)  null,
        COLUMN_55          varchar(24) null,
        COLUMN_56         varchar(5)  null,
        COLUMN_57          varchar(24) null,
        COLUMN_58         varchar(5)  null,
        COLUMN_59          varchar(24) null,
        COLUMN_60         varchar(5)  null,
        COLUMN_61          varchar(24) null,
        COLUMN_62         varchar(5)  null,
        COLUMN_63          varchar(24) null,
        COLUMN_64         varchar(5)  null,
        COLUMN_65          varchar(4)  null,
        COLUMN_66          varchar(24) null,
        COLUMN_67          varchar(24) null,
        COLUMN_68          varchar(4)  null,
        COLUMN_69          varchar(4)  null,
        COLUMN_70          varchar(4)  null,
        COLUMN_71          varchar(24) null,
        COLUMN_72          varchar(24) null,
        primary key (CODE, SAVE_DATE, SAVE_TIME, SERIALNUMBER)';
    set @createTable = concat('create table if not exists ', tableName,'(',@fields,')');
    prepare createTable from @createTable;
    execute createTable;
    deallocate prepare createTable;

    -- Export Records File (Canceled)
#     set @export = concat('select * into outfile ''',filePath,''' fields terminated by '','' lines terminated by ''\\n'' from ',tableName);
#     prepare exportFile from @export;
#     execute exportFile;
#     deallocate prepare exportFile;
end;
//
delimiter ;


drop procedure GenerateTableAndGetRecord;
