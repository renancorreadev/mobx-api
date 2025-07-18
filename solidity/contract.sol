// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

/**
 * @title VFinanceRegistry
 * @dev Sistema simplificado de registro de contratos de financiamento
 * @notice Dados essenciais on-chain, metadados servidos via API usando hash como identificador
 */
contract VFinanceRegistry {
    
    // =============================================================
    //                      STRUCTS
    // =============================================================
    
    struct ContractRecord {
        string regConId;           // ID único do contrato
        string numeroContrato;     // Número do contrato
        string dataContrato;       // Data do contrato
        bytes32 metadataHash;      // Hash único para buscar metadados via API
        uint256 timestamp;         // Timestamp do registro
        address registeredBy;      // Endereço que registrou
        bool active;               // Status ativo/inativo
    }
    
    // =============================================================
    //                      STATE VARIABLES
    // =============================================================
    
    mapping(string => ContractRecord) public contracts;
    mapping(string => bool) public contractExists;
    mapping(bytes32 => string) public hashToRegConId;
    mapping(string => bytes32) public regConIdToHash;
    string[] public contractIds;
    uint256 public totalContracts;
    
    // Servidor de API autorizado para atualizações (agora é URL string)
    string public apiServerUrl;
    address public apiServerAddress;  // Endereço autorizado para chamadas
    address public owner;
    
    // =============================================================
    //                         EVENTS
    // =============================================================
    
    event ContractRegistered(
        string indexed regConId,
        string indexed numeroContrato,
        bytes32 indexed metadataHash,
        address registeredBy,
        uint256 timestamp
    );
    
    event MetadataHashUpdated(
        string indexed regConId,
        bytes32 oldHash,
        bytes32 newHash,
        uint256 timestamp
    );
    
    event ContractStatusUpdated(
        string indexed regConId,
        bool active,
        uint256 timestamp
    );
    
    event ApiServerUpdated(
        string oldUrl,
        string newUrl,
        address oldAddress,
        address newAddress
    );
    
    // =============================================================
    //                         ERRORS
    // =============================================================
    
    error ContractAlreadyExists(string regConId);
    error ContractNotFound(string regConId);
    error InvalidInput(string paramName);
    error UnauthorizedAccess();
    error HashAlreadyUsed(bytes32 hash);
    
    // =============================================================
    //                        MODIFIERS
    // =============================================================
    
    modifier onlyOwner() {
        require(msg.sender == owner, "Not owner");
        _;
    }
    
    modifier onlyOwnerOrApiServer() {
        if (msg.sender != owner && msg.sender != apiServerAddress) {
            revert UnauthorizedAccess();
        }
        _;
    }
    
    // =============================================================
    //                        CONSTRUCTOR
    // =============================================================
    
    constructor(address initialOwner, string memory _apiServerUrl, address _apiServerAddress) {
        owner = initialOwner;
        apiServerUrl = _apiServerUrl;
        apiServerAddress = _apiServerAddress;
    }
    
    // =============================================================
    //                      METADATA URL FUNCTIONS
    // =============================================================
    
    /**
     * @dev Retorna a URL completa dos metadados para um hash específico
     * @param metadataHash Hash dos metadados
     * @return URL completa dos metadados
     */
    function tokenURI(bytes32 metadataHash) public view returns (string memory) {
        if (bytes(hashToRegConId[metadataHash]).length == 0) {
            revert InvalidInput("metadataHash");
        }
        
        return string(abi.encodePacked(
            apiServerUrl,
            "/api/metadata/0x",
            _toHexString(metadataHash)
        ));
    }
    
    /**
     * @dev Retorna a URL completa dos metadados para um regConId específico
     * @param regConId ID do contrato
     * @return URL completa dos metadados
     */
    function getMetadataUrl(string calldata regConId) public view returns (string memory) {
        if (!contractExists[regConId]) {
            revert ContractNotFound(regConId);
        }
        
        bytes32 metadataHash = regConIdToHash[regConId];
        return tokenURI(metadataHash);
    }
    
    /**
     * @dev Converte bytes32 para string hexadecimal
     * @param hash Hash para converter
     * @return String hexadecimal
     */
    function _toHexString(bytes32 hash) internal pure returns (string memory) {
        bytes memory buffer = new bytes(64);
        for (uint256 i = 0; i < 32; i++) {
            buffer[i * 2] = _toHexChar(uint8(hash[i]) / 16);
            buffer[i * 2 + 1] = _toHexChar(uint8(hash[i]) % 16);
        }
        return string(buffer);
    }
    
    /**
     * @dev Converte um dígito para caractere hexadecimal
     * @param digit Dígito para converter
     * @return Caractere hexadecimal
     */
    function _toHexChar(uint8 digit) internal pure returns (bytes1) {
        if (digit < 10) {
            return bytes1(uint8(bytes1('0')) + digit);
        }
        return bytes1(uint8(bytes1('a')) + digit - 10);
    }
    
    // =============================================================
    //                      MAIN FUNCTIONS
    // =============================================================
    
    /**
     * @dev Registra um novo contrato e gera hash automaticamente
     * @param regConId ID único do contrato
     * @param numeroContrato Número do contrato
     * @param dataContrato Data do contrato
     * @return metadataHash Hash gerado automaticamente para identificar metadados na API
     */
    function registerContract(
        string calldata regConId,
        string calldata numeroContrato,
        string calldata dataContrato
    ) external onlyOwnerOrApiServer returns (bytes32 metadataHash) {
        if (contractExists[regConId]) {
            revert ContractAlreadyExists(regConId);
        }
        if (bytes(regConId).length == 0) {
            revert InvalidInput("regConId");
        }
        if (bytes(numeroContrato).length == 0) {
            revert InvalidInput("numeroContrato");
        }
        if (bytes(dataContrato).length == 0) {
            revert InvalidInput("dataContrato");
        }
        
        // Gera hash único usando keccak256 com dados do contrato + timestamp + sender
        metadataHash = keccak256(
            abi.encodePacked(
                regConId,
                numeroContrato,
                dataContrato,
                block.timestamp,
                msg.sender,
                block.number
            )
        );
        
        // Verifica se o hash gerado já existe (muito improvável, mas por segurança)
        if (bytes(hashToRegConId[metadataHash]).length != 0) {
            // Se existir, adiciona um nonce para garantir unicidade
            metadataHash = keccak256(
                abi.encodePacked(
                    metadataHash,
                    totalContracts,
                    gasleft()
                )
            );
        }
        
        contracts[regConId] = ContractRecord({
            regConId: regConId,
            numeroContrato: numeroContrato,
            dataContrato: dataContrato,
            metadataHash: metadataHash,
            timestamp: block.timestamp,
            registeredBy: msg.sender,
            active: true
        });
        
        contractExists[regConId] = true;
        hashToRegConId[metadataHash] = regConId;
        regConIdToHash[regConId] = metadataHash;
        contractIds.push(regConId);
        totalContracts++;
        
        emit ContractRegistered(
            regConId,
            numeroContrato,
            metadataHash,
            msg.sender,
            block.timestamp
        );
        
        return metadataHash;
    }
    
    /**
     * @dev Atualiza o hash dos metadados
     * @param regConId ID do contrato
     * @param newMetadataHash Novo hash dos metadados
     */
    function updateMetadataHash(
        string calldata regConId,
        bytes32 newMetadataHash
    ) external onlyOwnerOrApiServer {
        if (!contractExists[regConId]) {
            revert ContractNotFound(regConId);
        }
        if (newMetadataHash == bytes32(0)) {
            revert InvalidInput("newMetadataHash");
        }
        if (bytes(hashToRegConId[newMetadataHash]).length != 0) {
            revert HashAlreadyUsed(newMetadataHash);
        }
        
        bytes32 oldHash = contracts[regConId].metadataHash;
        
        // Remove o mapeamento antigo
        delete hashToRegConId[oldHash];
        
        // Atualiza com o novo hash
        contracts[regConId].metadataHash = newMetadataHash;
        hashToRegConId[newMetadataHash] = regConId;
        regConIdToHash[regConId] = newMetadataHash;
        
        emit MetadataHashUpdated(regConId, oldHash, newMetadataHash, block.timestamp);
    }
    
    /**
     * @dev Atualiza o status do contrato
     * @param regConId ID do contrato
     * @param active Novo status
     */
    function updateContractStatus(
        string calldata regConId,
        bool active
    ) external onlyOwnerOrApiServer {
        if (!contractExists[regConId]) {
            revert ContractNotFound(regConId);
        }
        
        contracts[regConId].active = active;
        
        emit ContractStatusUpdated(regConId, active, block.timestamp);
    }
    
    /**
     * @dev Atualiza o servidor de API
     * @param newApiServerUrl Nova URL do servidor
     * @param newApiServerAddress Novo endereço autorizado
     */
    function updateApiServer(string memory newApiServerUrl, address newApiServerAddress) external onlyOwner {
        string memory oldUrl = apiServerUrl;
        address oldAddress = apiServerAddress;
        
        apiServerUrl = newApiServerUrl;
        apiServerAddress = newApiServerAddress;
        
        emit ApiServerUpdated(oldUrl, newApiServerUrl, oldAddress, newApiServerAddress);
    }

      
    /**
     * @dev Atualiza metadados e gera novo hash automaticamente
     * @param regConId ID do contrato
     * @param newDataContrato Novos dados do contrato
     * @return newMetadataHash Novo hash gerado automaticamente
     */
    function updateMetadataWithNewHash(
        string calldata regConId,
        string calldata newDataContrato
    ) external onlyOwnerOrApiServer returns (bytes32 newMetadataHash) {
        if (!contractExists[regConId]) {
            revert ContractNotFound(regConId);
        }
        if (bytes(newDataContrato).length == 0) {
            revert InvalidInput("newDataContrato");
        }
        
        ContractRecord storage contractRecord = contracts[regConId];
        bytes32 oldHash = contractRecord.metadataHash;
        
        // Gera novo hash único usando keccak256 com novos dados + timestamp + sender
        newMetadataHash = keccak256(
            abi.encodePacked(
                regConId,
                contractRecord.numeroContrato,
                newDataContrato,
                block.timestamp,
                msg.sender,
                block.number
            )
        );
        
        // Verifica se o hash gerado já existe (muito improvável, mas por segurança)
        if (bytes(hashToRegConId[newMetadataHash]).length != 0) {
            // Se existir, adiciona um nonce para garantir unicidade
            newMetadataHash = keccak256(
                abi.encodePacked(
                    newMetadataHash,
                    totalContracts,
                    gasleft()
                )
            );
        }
        
        // Remove o mapeamento antigo
        delete hashToRegConId[oldHash];
        
        // Atualiza com os novos dados e hash
        contractRecord.dataContrato = newDataContrato;
        contractRecord.metadataHash = newMetadataHash;
        hashToRegConId[newMetadataHash] = regConId;
        regConIdToHash[regConId] = newMetadataHash;
        
        emit MetadataHashUpdated(regConId, oldHash, newMetadataHash, block.timestamp);
        
        return newMetadataHash;
    }

    
    // =============================================================
    //                      VIEW FUNCTIONS
    // =============================================================
    
    /**
     * @dev Busca contrato por ID
     * @param regConId ID do contrato
     * @return Dados do contrato
     */
    function getContract(string calldata regConId)
        external
        view
        returns (ContractRecord memory)
    {
        if (!contractExists[regConId]) {
            revert ContractNotFound(regConId);
        }
        return contracts[regConId];
    }
    
    /**
     * @dev Busca contrato por hash de metadados
     * @param metadataHash Hash dos metadados
     * @return Dados do contrato
     */
    function getContractByHash(bytes32 metadataHash)
        external
        view
        returns (ContractRecord memory)
    {
        string memory regConId = hashToRegConId[metadataHash];
        if (bytes(regConId).length == 0) {
            revert ContractNotFound("hash not found");
        }
        return contracts[regConId];
    }
    
    /**
     * @dev Busca hash por ID do contrato
     * @param regConId ID do contrato
     * @return Hash dos metadados
     */
    function getHashByRegConId(string calldata regConId)
        external
        view
        returns (bytes32)
    {
        if (!contractExists[regConId]) {
            revert ContractNotFound(regConId);
        }
        return regConIdToHash[regConId];
    }
    
    /**
     * @dev Verifica se contrato existe
     * @param regConId ID do contrato
     * @return True se existe
     */
    function doesContractExist(string calldata regConId)
        external
        view
        returns (bool)
    {
        return contractExists[regConId];
    }
    
    /**
     * @dev Verifica se hash existe
     * @param metadataHash Hash dos metadados
     * @return True se existe
     */
    function doesHashExist(bytes32 metadataHash)
        external
        view
        returns (bool)
    {
        return bytes(hashToRegConId[metadataHash]).length != 0;
    }
    
    /**
     * @dev Retorna total de contratos
     * @return Número total de contratos
     */
    function getTotalContracts() external view returns (uint256) {
        return totalContracts;
    }
    
    /**
     * @dev Busca ID do contrato por índice
     * @param index Índice no array
     * @return ID do contrato
     */
    function getContractIdByIndex(uint256 index)
        external
        view
        returns (string memory)
    {
        if (index >= contractIds.length) {
            revert InvalidInput("index");
        }
        return contractIds[index];
    }
    
    /**
     * @dev Retorna contratos ativos paginados
     * @param offset Offset para paginação
     * @param limit Limite de resultados
     * @return Array de IDs de contratos ativos
     */
    function getActiveContracts(uint256 offset, uint256 limit)
        external
        view
        returns (string[] memory)
    {
        if (limit == 0 || limit > 100) {
            revert InvalidInput("limit");
        }
        
        uint256 activeCount = 0;
        
        // Conta contratos ativos
        for (uint256 i = 0; i < contractIds.length; i++) {
            if (contracts[contractIds[i]].active) {
                activeCount++;
            }
        }
        
        if (offset >= activeCount) {
            return new string[](0);
        }
        
        uint256 resultSize = activeCount - offset;
        if (resultSize > limit) {
            resultSize = limit;
        }
        
        string[] memory result = new string[](resultSize);
        uint256 currentIndex = 0;
        uint256 resultIndex = 0;
        
        for (uint256 i = 0; i < contractIds.length && resultIndex < resultSize; i++) {
            if (contracts[contractIds[i]].active) {
                if (currentIndex >= offset) {
                    result[resultIndex] = contractIds[i];
                    resultIndex++;
                }
                currentIndex++;
            }
        }
        
        return result;
    }

}
  