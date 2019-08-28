from flask import Blueprint, jsonify, request, url_for
from models import Sound
from app_oauth import require_oauth
import json
from models import licences as track_licenses


bp_api_v1_timelines = Blueprint("bp_api_v1_timelines", __name__)


@bp_api_v1_timelines.route("/api/v1/timelines/home", methods=["GET"])
@require_oauth(None)
def home():
    """
    Statuses from accounts the user follows.
    ---
    tags:
        - Timelines
    parameters:
        - name: count
          in: query
          type: integer
          required: true
          description: count
        - name: with_muted
          in: query
          type: boolean
          required: true
          description: with muted users
        - name: since_id
          in: query
          type: string
          description: last ID
    responses:
        200:
            description: Returns array of Status
    """
    count = int(request.args.get("count"), 20)
    since_id = request.args.get("since_id", None)

    # Get logged in user from bearer token, or None if not logged in
    # current_user = current_token.user
    sounds = Sound.query.filter(Sound.private.is_(False), Sound.transcode_state == Sound.TRANSCODE_DONE).order_by(
        Sound.uploaded
    )
    if since_id:
        sounds = sounds.filter(Sound.flake_id >= since_id)

    resp = []

    for sound in sounds.limit(count):
        si = sound.sound_infos.first()

        url_orig = url_for("get_uploads_stuff", thing="sounds", stuff=sound.path_sound(orig=True))
        url_transcode = url_for("get_uploads_stuff", thing="sounds", stuff=sound.path_sound(orig=False))

        track_obj = {
            "id": sound.flake_id,
            "title": sound.title,
            "user": sound.user.name,
            "description": sound.description,
            "picture_url": None,  # FIXME not implemented yet
            "media_orig": url_orig,
            "media_transcoded": url_transcode,
            "waveform": (json.loads(si.waveform) if si else None),
            "private": sound.private,
            "uploaded_on": sound.uploaded,
            "uploaded_elapsed": sound.elapsed(),
            "album_id": (sound.album.flake_id if sound.album else None),
            "processing": {
                "basic": (si.done_basic if si else None),
                "transcode_state": sound.transcode_state,
                "transcode_needed": sound.transcode_needed,
                "done": sound.processing_done(),
            },
            "metadatas": {
                "licence": track_licenses[sound.licence],
                "duration": (si.duration if si else None),
                "type": (si.type if si else None),
                "codec": (si.codec if si else None),
                "format": (si.format if si else None),
                "channels": (si.channels if si else None),
                "rate": (si.rate if si else None),  # Hz
            },
        }
        if si:
            if si.bitrate and si.bitrate_mode:
                track_obj["metadatas"]["bitrate"] = si.bitrate
                track_obj["metadatas"]["bitrate_mode"] = si.bitrate_mode

        resp.append(track_obj)

    response = jsonify(resp)
    response.mimetype = "application/json; charset=utf-8"
    response.status_code = 200
    return response