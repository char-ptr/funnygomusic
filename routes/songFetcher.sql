select * from (select iso.file as f, iso.title, iso.length as duration, ial.name as album, ial.id as alid , ia.name as artist, iso.track
        from indexed_songs iso
                 left join indexed_albums ial on iso.album = ial.id
                 left join indexed_artists ia on iso.artist = ia.id
        where iso.id in @ids
        union all
        select iso.file as f, iso.title, iso.length as duration, ial.name as album, ial.id as alid, ia.name as artist, iso.track
        from indexed_albums ial
                 left join indexed_songs iso on iso.album = ial.id
                 left join indexed_artists ia on iso.artist = ia.id
        where ial.id in @ids
        union all
        select iso.file as f, iso.title, iso.length as duration, ial.name as album, ial.id as alid, ia.name as artist, iso.track
        from indexed_artists ia
                 left join indexed_songs iso on iso.artist = ia.id
                 left join indexed_albums ial on iso.album = ial.id
        where ia.id in @ids) as t order by t.alid, t.track