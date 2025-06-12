`timescale 1ns / 1ps
// VTX1 Cache Controller - Functional Implementation

// Include VTX1 common infrastructure
`include "vtx1_interfaces.v"
`include "vtx1_state_constants.v"
`include "vtx1_error_macros.v"

module cache_controller (
    input  wire                     clk,
    input  wire                     rst_n,
    
    // CPU Interface - Instruction Cache
    input  wire                     icache_req,
    input  wire [36-1:0] icache_addr,
    output reg  [36-1:0] icache_data,
    output reg                      icache_hit,
    output reg                      icache_ready,
    
    // CPU Interface - Data Cache
    input  wire                     dcache_req,
    input  wire                     dcache_wr,
    input  wire [1:0]               dcache_size,
    input  wire [36-1:0] dcache_addr,
    input  wire [36-1:0] dcache_wdata,
    output reg  [36-1:0] dcache_rdata,
    output reg                      dcache_hit,
    output reg                      dcache_ready,
    
    // Memory Controller Interface - for cache line fills
    output reg                      cache_req,
    output reg                      cache_wr,
    output reg  [36-1:0] cache_addr,
    output reg  [288-1:0] cache_wdata,
    input  wire [288-1:0] cache_rdata,
    input  wire                     cache_ready,
    input  wire                     cache_error,
    
    // Control and Status
    input  wire                     cache_enable,
    input  wire                     cache_flush,
    output reg  [3:0]               cache_state,
    
    // Debug and Monitoring (simplified for Phase 4)
    output reg  [31:0]              icache_hits,
    output reg  [31:0]              icache_misses,
    output reg  [31:0]              dcache_hits,
    output reg  [31:0]              dcache_misses
);

    // Cache Controller with Simplified but Functional Implementation
    
    // Cache storage - very basic direct-mapped cache
    // Instruction Cache: 64 entries of 32-bit words
    localparam ICACHE_SIZE = 64;
    localparam ICACHE_INDEX_BITS = 6;
    reg [36-1:0] icache_data_array [0:ICACHE_SIZE-1];
    reg [36-ICACHE_INDEX_BITS-2:0] icache_tag_array [0:ICACHE_SIZE-1];
    reg icache_valid_array [0:ICACHE_SIZE-1];
    
    // Data Cache: 64 entries of 32-bit words
    localparam DCACHE_SIZE = 64;
    localparam DCACHE_INDEX_BITS = 6;
    reg [36-1:0] dcache_data_array [0:DCACHE_SIZE-1];
    reg [36-DCACHE_INDEX_BITS-2:0] dcache_tag_array [0:DCACHE_SIZE-1];
    reg dcache_valid_array [0:DCACHE_SIZE-1];
    reg dcache_dirty_array [0:DCACHE_SIZE-1];
    
    // Internal state
    reg [3:0] cache_state_reg;
    reg [31:0] icache_hits_reg;
    reg [31:0] icache_misses_reg;
    reg [31:0] dcache_hits_reg;
    reg [31:0] dcache_misses_reg;
    reg filling_icache;
    reg filling_dcache;
    reg [36-1:0] fill_addr;
    reg [ICACHE_INDEX_BITS-1:0] fill_icache_index;
    reg [DCACHE_INDEX_BITS-1:0] fill_dcache_index;
    
    // Loop variable for initialization
    integer i;
    integer j;
    
    // Cache line fill state
    // reg filling_icache;
    // reg filling_dcache;
    // reg [36-1:0] fill_addr;
    // reg [ICACHE_INDEX_BITS-1:0] fill_icache_index;
    // reg [DCACHE_INDEX_BITS-1:0] fill_dcache_index;
    
    // Cache lookup logic
    wire [ICACHE_INDEX_BITS-1:0] icache_index;
    wire [36-ICACHE_INDEX_BITS-2:0] icache_tag;
    wire icache_tag_match;
    wire icache_cache_hit;
    
    wire [DCACHE_INDEX_BITS-1:0] dcache_index;
    wire [36-DCACHE_INDEX_BITS-2:0] dcache_tag;
    wire dcache_tag_match;
    wire dcache_cache_hit;
    
    // Address parsing
    assign icache_index = icache_addr[ICACHE_INDEX_BITS+1:2];
    assign icache_tag = icache_addr[36-1:ICACHE_INDEX_BITS+2];
    assign icache_tag_match = (icache_tag_array[icache_index] == icache_tag);
    assign icache_cache_hit = icache_req && icache_valid_array[icache_index] && icache_tag_match;
    
    assign dcache_index = dcache_addr[DCACHE_INDEX_BITS+1:2];
    assign dcache_tag = dcache_addr[36-1:DCACHE_INDEX_BITS+2];
    assign dcache_tag_match = (dcache_tag_array[dcache_index] == dcache_tag);
    assign dcache_cache_hit = dcache_req && dcache_valid_array[dcache_index] && dcache_tag_match;
    
    // State machine
    always @(posedge clk or negedge rst_n) begin
        if (!rst_n) begin
            cache_state_reg <= 4'b0000;
            icache_hits_reg <= 32'b0;
            icache_misses_reg <= 32'b0;
            dcache_hits_reg <= 32'b0;
            dcache_misses_reg <= 32'b0;
            filling_icache <= 1'b0;
            filling_dcache <= 1'b0;
            fill_addr <= {36{1'b0}};
            fill_icache_index <= {ICACHE_INDEX_BITS{1'b0}};
            fill_dcache_index <= {DCACHE_INDEX_BITS{1'b0}};
            
            // Initialize cache arrays to invalid
            for (i = 0; i < ICACHE_SIZE; i = i + 1) begin
                icache_valid_array[i] <= 1'b0;
                icache_data_array[i] <= 36'b010101010101010101010101010101010101;
                icache_tag_array[i] <= {(36-ICACHE_INDEX_BITS-2){1'b0}};
            end
            for (i = 0; i < DCACHE_SIZE; i = i + 1) begin
                dcache_valid_array[i] <= 1'b0;
                dcache_dirty_array[i] <= 1'b0;
                dcache_data_array[i] <= 36'b010101010101010101010101010101010101;
                dcache_tag_array[i] <= {(36-DCACHE_INDEX_BITS-2){1'b0}};
            end
        end else begin
            case (cache_state_reg)
                4'b0000: begin
                    if (cache_enable) begin
                        cache_state_reg <= 4'b0001;
                    end
                end
                
                4'b0001: begin
                    if (!cache_enable) begin
                        cache_state_reg <= 4'b0000;
                    end else if (cache_flush) begin
                        // Handle cache flush
                        for (j = 0; j < ICACHE_SIZE; j = j + 1) begin
                            icache_valid_array[j] <= 1'b0;
                        end
                        for (j = 0; j < DCACHE_SIZE; j = j + 1) begin
                            dcache_valid_array[j] <= 1'b0;
                            dcache_dirty_array[j] <= 1'b0;
                        end
                        icache_hits_reg <= 32'b0;
                        icache_misses_reg <= 32'b0;
                        dcache_hits_reg <= 32'b0;
                        dcache_misses_reg <= 32'b0;
                    end else begin
                        // Handle cache requests
                        if (icache_req && !icache_cache_hit && !filling_icache) begin
                            // Instruction cache miss - start fill
                            cache_state_reg <= 4'b0010;
                            filling_icache <= 1'b1;
                            fill_addr <= {icache_addr[36-1:2], 2'b00}; // Word aligned
                            fill_icache_index <= icache_index;
                            icache_misses_reg <= icache_misses_reg + 1;
                        end else if (dcache_req && !dcache_cache_hit && !filling_dcache) begin
                            // Data cache miss - start fill
                            cache_state_reg <= 4'b0010;
                            filling_dcache <= 1'b1;
                            fill_addr <= {dcache_addr[36-1:2], 2'b00}; // Word aligned
                            fill_dcache_index <= dcache_index;
                            dcache_misses_reg <= dcache_misses_reg + 1;
                        end else begin
                            // Handle hits
                            if (icache_cache_hit) begin
                                icache_hits_reg <= icache_hits_reg + 1;
                            end
                            if (dcache_cache_hit) begin
                                dcache_hits_reg <= dcache_hits_reg + 1;
                                // Handle data cache write
                                if (dcache_wr) begin
                                    dcache_data_array[dcache_index] <= dcache_wdata;
                                    dcache_dirty_array[dcache_index] <= 1'b1;
                                end
                            end
                        end
                    end
                end
                
                4'b0010: begin
                    if (cache_ready && !cache_error) begin
                        // Fill completed successfully
                        if (filling_icache) begin
                            icache_data_array[fill_icache_index] <= cache_rdata[36-1:0];
                            icache_tag_array[fill_icache_index] <= fill_addr[36-1:ICACHE_INDEX_BITS+2];
                            icache_valid_array[fill_icache_index] <= 1'b1;
                            filling_icache <= 1'b0;
                        end
                        if (filling_dcache) begin
                            dcache_data_array[fill_dcache_index] <= cache_rdata[36-1:0];
                            dcache_tag_array[fill_dcache_index] <= fill_addr[36-1:DCACHE_INDEX_BITS+2];
                            dcache_valid_array[fill_dcache_index] <= 1'b1;
                            dcache_dirty_array[fill_dcache_index] <= 1'b0;
                            filling_dcache <= 1'b0;
                        end
                        cache_state_reg <= 4'b0001;
                    end else if (cache_error) begin
                        // Fill failed - return to ready and let subsequent requests retry
                        filling_icache <= 1'b0;
                        filling_dcache <= 1'b0;
                        cache_state_reg <= 4'b0001;
                    end
                end
                
                default: begin
                    cache_state_reg <= 4'b0000;
                end
            endcase
        end
    end
    
    // Combinational outputs
    always @(*) begin
        // Default values
        icache_data = 36'b010101010101010101010101010101010101;
        icache_hit = 1'b0;
        icache_ready = 1'b0;
        dcache_rdata = 36'b010101010101010101010101010101010101;
        dcache_hit = 1'b0;
        dcache_ready = 1'b0;
        cache_req = 1'b0;
        cache_wr = 1'b0;
        cache_addr = {36{1'b0}};
        cache_wdata = {288{1'b0}};

        case (cache_state_reg)
            4'b0001: begin
                // Handle instruction cache
                if (icache_req) begin
                    if (icache_cache_hit) begin
                        icache_data = icache_data_array[icache_index];
                        icache_hit = 1'b1;
                        icache_ready = 1'b1;
                    end else if (!filling_icache) begin
                        // Cache miss - ready is low until fill completes
                        icache_ready = 1'b0;
                    end
                end
                
                // Handle data cache
                if (dcache_req) begin
                    if (dcache_cache_hit) begin
                        dcache_rdata = dcache_data_array[dcache_index];
                        dcache_hit = 1'b1;
                        dcache_ready = 1'b1;
                    end else if (!filling_dcache) begin
                        // Cache miss - ready is low until fill completes
                        dcache_ready = 1'b0;
                    end
                end
            end
            
            4'b0010: begin
                // Request cache line fill from memory
                cache_req = 1'b1;
                cache_wr = 1'b0;
                cache_addr = fill_addr;
                // For reads, wdata doesn't matter
                cache_wdata = {288{1'b0}};
            end
        endcase
    end
    
    // Output assignments
    always @(*) begin
        cache_state = cache_state_reg;
        icache_hits = icache_hits_reg;
        icache_misses = icache_misses_reg;
        dcache_hits = dcache_hits_reg;
        dcache_misses = dcache_misses_reg;
    end

endmodule
