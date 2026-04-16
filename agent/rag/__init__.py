"""
Multi-channel retrieval engine for MiniBlog.

Inspired by CAgent's MultiChannelRetrievalEngine:
- Multiple search channels execute in parallel (asyncio.gather)
- Post-processors form a sequential pipeline (dedup -> RRF -> rerank -> filter)
- Strategy pattern: add new channels/processors by implementing the interface
"""
