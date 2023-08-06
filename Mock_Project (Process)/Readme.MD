**Procedure for DB**

create procedure GenerateTableAndGetRecord(
    in tableName varchar(255),
    in filePath longtext
)
begin
    -- Content for Procedure 'GenerateTableAndGetRecord()' in here
    set@fields= 'QCD            varchar(42) not null,
        TIME           varchar(5)  not null,
        TKZXD          date        not null,
        TKTIM          varchar(15) not null,
        TKSERIALNUMBER int         not null,
        TKQKBN         varchar(2)  not null,
        SNDC           varchar(5)  not null,
        ZXD            varchar(10) null,
        NOPV           varchar(24) null,
        ZSS            varchar(4)  null,
        DPS            varchar(4)  null,
        DPS_T          varchar(5)  null,
        DPP            varchar(24) null,
        DPP_T          varchar(5)  null,
        DPR            varchar(24) null,
        DPR_T          varchar(5)  null,
        XV             varchar(24) null,
        XV_T           varchar(5)  null,
        NOQ            varchar(24) null,
        QAS            varchar(4)  null,
        QAP            varchar(24) null,
        QAP_T          varchar(5)  null,
        QAR            varchar(24) null,
        QAR_T          varchar(5)  null,
        AV             varchar(24) null,
        AV_T           varchar(5)  null,
        QBS            varchar(4)  null,
        QBP            varchar(24) null,
        QBP_T          varchar(5)  null,
        QBR            varchar(24) null,
        QBR_T          varchar(5)  null,
        BV             varchar(24) null,
        BV_T           varchar(5)  null,
        J_SNDC         varchar(5)  null,
        XJ             varchar(24) null,
        XJ_T           varchar(5)  null,
        DPCF           varchar(5)  null,
        JQCS           varchar(4)  null,
        JQJS           varchar(4)  null,
        DPCY           varchar(24) null,
        DPCY_T         varchar(5)  null,
        TLNM           varchar(24) null,
        XV_S           varchar(8)  null,
        AV_S           varchar(8)  null,
        BV_S           varchar(8)  null,
        DPP_S          varchar(8)  null,
        J_ZXD          varchar(10) null,
        MIDP           varchar(24) null,
        MIDP_T         varchar(5)  null,
        DYWP           varchar(24) null,
        DYWP_T         varchar(5)  null,
        DYWR           varchar(24) null,
        VWAP           varchar(24) null,
        VWAP_T         varchar(5)  null,
        ABV            varchar(24) null,
        ABV_T          varchar(5)  null,
        AAV            varchar(24) null,
        AAV_T          varchar(5)  null,
        QOV            varchar(24) null,
        QOV_T          varchar(5)  null,
        QUV            varchar(24) null,
        QUV_T          varchar(5)  null,
        INAV           varchar(24) null,
        INAV_T         varchar(5)  null,
        IYRP           varchar(24) null,
        IYRP_T         varchar(5)  null,
        IQRP           varchar(24) null,
        IQRP_T         varchar(5)  null,
        QACY           varchar(24) null,
        QACY_T         varchar(5)  null,
        QBCY           varchar(24) null,
        QBCY_T         varchar(5)  null,
        TSTF           varchar(4)  null,
        NO             varchar(24) null,
        `NOT`          varchar(24) null,
        1KF            varchar(4)  null,
        2KF            varchar(4)  null,
        DKF            varchar(4)  null,
        DSYRP          varchar(24) null,
        DSYWP          varchar(24) null,
        primary key (QCD, TKZXD, TKTIM, TKSERIALNUMBER)';
    set @createTable = concat('create table if not exists ', tableName,'(',@fields,')');
    prepare createTable from @createTable;
    execute createTable;
    deallocate prepare createTable;

    -- Export Records File
    set @export = concat('select * into outfile ''',filePath,''' fields terminated by '','' lines terminated by ''\\n'' from ',tableName);
    prepare exportFile from @export;
    execute exportFile;
    deallocate prepare exportFile;
end;
//
delimiter ;


drop procedure GenerateTableAndGetRecord;
