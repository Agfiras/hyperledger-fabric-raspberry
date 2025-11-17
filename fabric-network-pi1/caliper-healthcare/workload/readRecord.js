'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

class ReadRecordWorkload extends WorkloadModuleBase {
    constructor() {
        super();
        this.recordIndex = 0;
        this.workerIndex = 0;
        this.totalWorkers = 5;
        this.totalRounds = 30;
        this.maxRecordsPerWorker = 500;
    }

    async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
        await super.initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext);
        this.workerIndex = workerIndex;
        this.totalWorkers = totalWorkers;
        this.recordIndex = 0;
        
        if (roundArguments && roundArguments.maxRecords) {
            this.maxRecordsPerWorker = roundArguments.maxRecords;
        }
    }

    async submitTransaction() {
        this.recordIndex++;
        
        // Match CreateRecord format
        const randomRound = Math.floor(Math.random() * this.totalRounds);
        const randomWorker = Math.floor(Math.random() * this.totalWorkers);
        const randomRecord = Math.floor(Math.random() * this.maxRecordsPerWorker) + 1;
        
        const recordID = 'REC_R' + randomRound + '_W' + randomWorker + '_' + randomRecord.toString().padStart(6, '0');

        const request = {
            contractId: 'healthcc',
            contractFunction: 'ReadRecord',
            invokerIdentity: 'User1',
            contractArguments: [recordID],   // ✅ Fixed: use recordID
            readOnly: true
        };

        await this.sutAdapter.sendRequests(request);
    }
}

function createWorkloadModule() {
    return new ReadRecordWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;