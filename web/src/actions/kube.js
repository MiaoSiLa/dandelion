import { RSAA } from 'redux-api-middleware'
import querystring from 'querystring'

import { API_URL, apiPayload } from '../api'

export const KUBE_LIST_REQUEST = 'KUBE_LIST_REQUEST'
export const KUBE_LIST_SUCCESS = 'KUBE_LIST_SUCCESS'
export const KUBE_LIST_FAILURE = 'KUBE_LIST_FAILURE'

export function kubeListDeployments() {
  return {
    [RSAA]: {
      method: 'GET',
      endpoint: API_URL + '/kube/list',
      types: [
        KUBE_LIST_REQUEST,
        { type: KUBE_LIST_SUCCESS, payload: apiPayload },
        KUBE_LIST_FAILURE
      ],
      credentials: 'include',
    }
  }
}

export const KUBE_LIST_TAGS_REQUEST = 'KUBE_LIST_TAGS_REQUEST'
export const KUBE_LIST_TAGS_SUCCESS = 'KUBE_LIST_TAGS_SUCCESS'
export const KUBE_LIST_TAGS_FAILURE = 'KUBE_LIST_TAGS_FAILURE'

export function kubeListTags(name) {
  return {
    [RSAA]: {
      method: 'GET',
      endpoint: API_URL + '/kube/listtags/' + name,
      types: [
        KUBE_LIST_TAGS_REQUEST,
        { type: KUBE_LIST_TAGS_SUCCESS, payload: apiPayload, meta: { name } },
        KUBE_LIST_TAGS_FAILURE
      ],
      credentials: 'include',
    }
  }
}

export const KUBE_SET_TAG_REQUEST = 'KUBE_SET_TAG_REQUEST'
export const KUBE_SET_TAG_SUCCESS = 'KUBE_SET_TAG_SUCCESS'
export const KUBE_SET_TAG_FAILURE = 'KUBE_SET_TAG_FAILURE'

export function kubeSetTag(name, tag) {
  return {
    [RSAA]: {
      method: 'POST',
      endpoint: API_URL + '/kube/setversiontag/' + name,
      types: [
        KUBE_SET_TAG_REQUEST,
        { type: KUBE_SET_TAG_SUCCESS, payload: apiPayload, meta: { name, tag } },
        KUBE_SET_TAG_FAILURE
      ],
      credentials: 'include',
      headers: {
        'Content-Type' : 'application/x-www-form-urlencoded',
      },
      body: querystring.stringify({
        version_tag: tag,
      }),
    }
  }
}

export const KUBE_ROLLBACK_REQUEST = 'KUBE_ROLLBACK_REQUEST'
export const KUBE_ROLLBACK_SUCCESS = 'KUBE_ROLLBACK_SUCCESS'
export const KUBE_ROLLBACK_FAILURE = 'KUBE_ROLLBACK_FAILURE'

export function kubeRollback(name) {
  return {
    [RSAA]: {
      method: 'POST',
      endpoint: API_URL + '/kube/rollback/' + name,
      types: [
        KUBE_ROLLBACK_REQUEST,
        { type: KUBE_ROLLBACK_SUCCESS, payload: apiPayload, meta: { name } },
        KUBE_ROLLBACK_FAILURE
      ],
      credentials: 'include',
    }
  }
}
