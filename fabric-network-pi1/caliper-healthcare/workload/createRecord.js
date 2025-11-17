'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

class CreateRecordWorkload extends WorkloadModuleBase {
    constructor() {
        super();
        this.recordIndex = 0;
        this.workerIndex = 0;
        this.roundIndex = 0;
    }

    async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
        await super.initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext);
        this.workerIndex = workerIndex;
        this.roundIndex = roundIndex;
        this.recordIndex = 0;
    }

    async submitTransaction() {
        this.recordIndex++;
        
        // Generate globally unique ID for the key
        const recordID = 'REC_R' + this.roundIndex + '_W' + this.workerIndex + '_' + this.recordIndex.toString().padStart(6, '0');
        
        // Create an object for the container's data
        const containerData = {
            patientName: 'Patient_' + recordID,
            diagnosis: 'Diagnosis_' + (this.recordIndex % 10),
            treatment: 'Treatment_' + (this.recordIndex % 5),
            timestamp: new Date().toISOString()
        };

        // The second argument should be a string. We'll stringify the JSON object.
        const containerDataString = JSON.stringify(containerData);

        const request = {
            contractId: 'healthcc',
            contractFunction: 'CreateRecord',  // Changed from CreateOrUpdateContainer
            invokerIdentity: 'User1',
            contractArguments: [
                recordID,                              // param0: id
                containerData.patientName,             // param1: patientName
                containerData.diagnosis,               // param2: diagnosis
                containerData.treatment                // param3: treatment
    ],
    readOnly: false
};

        await this.sutAdapter.sendRequests(request);
    }
}

function createWorkloadModule() {
    return new CreateRecordWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;